package commands

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var downloadCommand = discord.SlashCommandCreate{
	Name:        "download",
	Description: "Download a video",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "url",
			Description: "The URL of the video",
			Required:    true,
		},
	},
}

type DownloadTask struct {
	e                 *handler.CommandEvent
	url               string
	filePathProcessed string
	tempDir           string
}

var taskQueue = make(chan DownloadTask, 10)
var once sync.Once

func HandleDownload(b *wokkibot.Wokkibot) handler.CommandHandler {
	once.Do(func() {
		go downloadWorker()
	})

	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		if data.String("url") == "" {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No URL provided").Build())
		}

		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		randomFileName := RandomName(10)
		processedFilename := fmt.Sprintf("%v_processed.mp4", randomFileName)

		if _, err := os.Stat("downloads"); os.IsNotExist(err) {
			if err := os.MkdirAll("downloads", 0755); err != nil {
				return fmt.Errorf("failed to create downloads directory: %w", err)
			}
		}

		tempDir, err := os.MkdirTemp("downloads", "video_*")
		if err != nil {
			fmt.Printf("Failed to create temp directory: %v\n", err)
			return err
		}

		filePathProcessed := filepath.Join(tempDir, processedFilename)

		task := DownloadTask{
			e:                 e,
			url:               data.String("url"),
			filePathProcessed: filePathProcessed,
			tempDir:           tempDir,
		}
		taskQueue <- task

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Waiting for previous download tasks to finish...").Build())
		return err
	}
}

func downloadWorker() {
	for task := range taskQueue {
		handleDownloadAndConversion(task)
	}
}

func handleDownloadAndConversion(task DownloadTask) {
	e := task.e

	cleanup := func() {
		if err := os.RemoveAll(task.tempDir); err != nil {
			fmt.Printf("Error while removing downloaded files: %v", err)
		}
	}

	defer cleanup()

	e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Downloading video...").Build())

	downloadCtx, downloadCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer downloadCancel()

	downloadOutput := filepath.Join(task.tempDir, "video_download.%(ext)s")
	cmd := exec.CommandContext(downloadCtx, "yt-dlp", task.url, "-o", downloadOutput)
	if err := cmd.Run(); err != nil {
		if downloadCtx.Err() == context.DeadlineExceeded {
			e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Download canceled as it took too long").Build())
		}
		e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while downloading video").Build())
		return
	}

	downloadedFiles, err := filepath.Glob(filepath.Join(task.tempDir, "video_download.*"))
	if err != nil {
		e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while finding downloaded file").Build())
		return
	}
	inputFilePath := downloadedFiles[0]

	conversionCtx, conversionCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer conversionCancel()

	e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Converting video...").Build())
	conversion := exec.CommandContext(conversionCtx, "ffmpeg",
		"-i", inputFilePath,
		"-c:v", "h264",
		"-b:v", "1M",
		"-c:a", "aac",
		"-pix_fmt", "yuv420p",
		"-f", "mp4",
		task.filePathProcessed,
		"-progress", "pipe:1", "-nostats",
	)

	totalDuration, _ := getVideoDuration(inputFilePath)

	conversionStdout, err := conversion.StdoutPipe()
	if err != nil {
		e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while getting stdout").Build())
		return
	}

	if err := conversion.Start(); err != nil {
		e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while starting conversion").Build())
		return
	}

	ffmpegScanner := bufio.NewScanner(conversionStdout)
	var currentTime string

	lastUpdateTime := time.Now()
	var lastPercentage float64

	for ffmpegScanner.Scan() {
		line := ffmpegScanner.Text()

		if strings.Contains(line, "out_time=") {
			timeIndex := strings.Index(line, "out_time=")
			if timeIndex != -1 {
				currentTime = line[timeIndex+9:]

				var progressTime float64
				parts := strings.Split(currentTime, ":")
				if len(parts) == 3 {
					hours, _ := strconv.ParseFloat(parts[0], 64)
					minutes, _ := strconv.ParseFloat(parts[1], 64)
					seconds, _ := strconv.ParseFloat(parts[2], 64)
					progressTime = hours*3600 + minutes*60 + seconds
				}

				if totalDuration > 0 {
					progressPercentage := (progressTime / totalDuration) * 100
					if time.Since(lastUpdateTime) >= 1*time.Second && progressPercentage != lastPercentage {
						progress := progressBar(progressPercentage)
						e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent(fmt.Sprintf("Converting video\n%s %.2f%%", progress, progressPercentage)).Build())
						lastUpdateTime = time.Now()
						lastPercentage = progressPercentage
					}
				}
			}
		}
	}

	if err := conversion.Wait(); err != nil {
		if conversionCtx.Err() == context.DeadlineExceeded {
			e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Conversion canceled as it took too long").Build())
		}
		e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while processing video").Build())
		return
	}

	file, err := os.Open(task.filePathProcessed)
	if err != nil {
		e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while opening file").Build())
		return
	}
	defer file.Close()

	_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("").AddFile(task.filePathProcessed, task.filePathProcessed, file).Build())
	if err != nil {
		e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while attaching file").Build())
	}
}

func RandomName(length int) string {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = characters[seed.Intn(len(characters))]
	}
	return string(b)
}

func getVideoDuration(videoFile string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoFile)

	durationOutput, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("error getting duration: %w", err)
	}

	totalDuration, err := strconv.ParseFloat(strings.TrimSpace(string(durationOutput)), 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing duration: %w", err)
	}

	return totalDuration, nil
}

func progressBar(percentage float64) string {
	filledBlocks := int(percentage / 100 * float64(20))

	bar := strings.Repeat("█", filledBlocks) + strings.Repeat("░", 20-filledBlocks)

	return bar
}
