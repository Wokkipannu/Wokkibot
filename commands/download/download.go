package download

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var DownloadCommand = discord.SlashCommandCreate{
	Name:        "download",
	Description: "Download a video",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "url",
			Description: "The URL of the video",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "resolution",
			Description: "Overwrite the default 720p resolution in format sort",
			Required:    false,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "1080p",
					Value: "1080",
				},
				{
					Name:  "720p",
					Value: "720",
				},
				{
					Name:  "480p",
					Value: "480",
				},
				{
					Name:  "360p",
					Value: "360",
				},
				{
					Name:  "240p",
					Value: "240",
				},
				{
					Name:  "144p",
					Value: "144",
				},
			},
		},
		discord.ApplicationCommandOptionString{
			Name:        "start",
			Description: "The time to start the video from (e.g. 1:30 or 90)",
			Required:    false,
		},
		discord.ApplicationCommandOptionString{
			Name:        "end",
			Description: "The time to end the video at (e.g. 2:45 or 165)",
			Required:    false,
		},
	},
}

type DownloadTask struct {
	e                 *handler.CommandEvent
	url               string
	filePathProcessed string
	tempDir           string
	maxFileSize       int
	resolution        string
	from              string
	to                string
}

type DownloadProgress struct {
	ProgressPercentage string `json:"progress_percentage"`
}

const (
	downloadTimeout   = 3 * time.Minute
	conversionTimeout = 5 * time.Minute
	updateInterval    = 1 * time.Second
	defaultBitrate    = "1M"
	defaultResolution = "720"
)

var taskQueue = make(chan DownloadTask, 10)
var once sync.Once

func HandleDownload(b *wokkibot.Wokkibot) handler.CommandHandler {
	once.Do(func() {
		go downloadWorker()
	})

	return func(e *handler.CommandEvent) error {
		url := e.SlashCommandInteractionData().String("url")

		if url == "" {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No URL provided").Build())
		}

		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		if err := os.MkdirAll("downloads", 0755); err != nil {
			return err
		}
		tempDir, err := os.MkdirTemp("downloads", "video_*")
		if err != nil {
			return err
		}

		guild, _ := e.Guild()

		newUrl, _ := handleSpecialScenarios(e, url)
		if newUrl != url {
			url = newUrl
		}

		var res string
		if e.SlashCommandInteractionData().String("resolution") != "" {
			res = e.SlashCommandInteractionData().String("resolution")
		} else {
			res = defaultResolution
		}

		task := DownloadTask{
			e:                 e,
			url:               url,
			filePathProcessed: filepath.Join(tempDir, fmt.Sprintf("%s_processed.mp4", utils.GenerateRandomName(10))),
			tempDir:           tempDir,
			maxFileSize:       utils.CalculateMaximumFileSizeForGuild(guild),
			resolution:        res,
			from:              e.SlashCommandInteractionData().String("start"),
			to:                e.SlashCommandInteractionData().String("end"),
		}
		taskQueue <- task

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent("Waiting for previous download tasks to finish...").
			Build())
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

	defer cleanup(task.tempDir)

	e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
		SetContent("Starting video download...").
		Build())

	var downloadedFile string
	var err error

	if strings.HasPrefix(task.url, "https://i.ylilauta.org/") {
		downloadedFile, err = downloadFileWithCurl(task, e)
	} else {
		downloadedFile, err = downloadFile(task, e)
	}

	if err != nil {
		utils.HandleError(e, "Error while downloading video", err.Error())
		return
	}

	processedFile, err := convertVideo(task, e, downloadedFile)
	if err != nil {
		utils.HandleError(e, "Error while converting video", err.Error())
		return
	}

	if err := attachFile(e, processedFile); err != nil {
		utils.HandleError(e, "Error while attaching file", err.Error())
		return
	}
}

func downloadFile(task DownloadTask, e *handler.CommandEvent) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	from := task.from
	to := task.to

	output := filepath.Join(task.tempDir, "video_download.%(ext)s")
	cmd := exec.CommandContext(ctx, "yt-dlp",
		task.url,
		"-o", output,
		"--max-filesize", fmt.Sprintf("%dM", task.maxFileSize),
		"--format-sort", fmt.Sprintf("res:%s,codec:h264", task.resolution),
		"--merge-output-format", "mp4",
		"--cookies", "cookies.txt",
		"--progress-template", "{\"progress_percentage\": \"%(progress._percent_str)s\"}",
		"--newline",
	)

	if from != "" {
		toValue := "inf"
		if to != "" {
			toValue = to
		}
		cmd.Args = append(cmd.Args, "--download-sections", fmt.Sprintf("*%s-%s", from, toValue))
	}

	return executeOperation(e, task, cmd, ctx, "download", "")
}

func downloadFileWithCurl(task DownloadTask, e *handler.CommandEvent) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	output := filepath.Join(task.tempDir, "video_download.%(ext)s")

	cmd := exec.CommandContext(ctx, "curl",
		"-L",
		"-f",
		"-#",
		"-o", output,
		"--max-filesize", fmt.Sprintf("%d", task.maxFileSize*1024*1024),
		task.url,
	)

	return executeOperation(e, task, cmd, ctx, "curldownload", "")
}

func convertVideo(task DownloadTask, e *handler.CommandEvent, downloadedFile string) (string, error) {
	codec, err := getCodec(downloadedFile)
	if err != nil {
		return "", err
	}

	if codec == "h264" {
		return downloadedFile, nil
	}

	if codec == "mp3" {
		return downloadedFile, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), conversionTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", downloadedFile,
		"-c:v", "h264",
		"-b:v", defaultBitrate,
		"-c:a", "aac",
		"-pix_fmt", "yuv420p",
		"-f", "mp4",
		task.filePathProcessed,
		"-progress", "pipe:1",
		"-nostats",
	)

	return executeOperation(e, task, cmd, ctx, "conversion", downloadedFile)
}

func executeOperation(e *handler.CommandEvent, task DownloadTask, cmd *exec.Cmd, ctx context.Context, operation string, downloadedFile string) (string, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("error getting stdout: %w", err)
	}

	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting command: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	lastUpdate := time.Now()
	var lastPercentage float64

	e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
		SetContent(fmt.Sprintf("Starting video %s\n%s %.2f%%", operation, createProgressBar(0.0), 0.0)).
		Build())

	for scanner.Scan() {
		if operation == "download" {
			if strings.Contains(scanner.Text(), "File is larger than max-filesize") {
				_ = cmd.Process.Kill()
				return "", fmt.Errorf("file size exceeds the maximum allowed size for this guild. Maximum is %dMB", task.maxFileSize)
			}

			if !strings.HasPrefix(scanner.Text(), "{") || !json.Valid([]byte(scanner.Text())) {
				continue
			}

			var progress DownloadProgress
			if err := json.Unmarshal([]byte(scanner.Text()), &progress); err != nil {
				continue
			}

			percentage, _ := strconv.ParseFloat(strings.TrimSuffix(progress.ProgressPercentage, "%"), 64)
			if time.Since(lastUpdate) >= updateInterval && percentage != lastPercentage {
				progress := createProgressBar(percentage)

				e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
					SetContent(fmt.Sprintf("Downloading video\n%s %.2f%%", progress, percentage)).
					Build())

				lastUpdate = time.Now()
				lastPercentage = percentage
			}
		}

		if operation == "conversion" {
			totalDuration, _ := getVideoDuration(downloadedFile)

			line := scanner.Text()

			if strings.Contains(line, "out_time=") {
				timeIndex := strings.Index(line, "out_time=")
				if timeIndex != -1 {
					currentTime := line[timeIndex+9:]

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
						if time.Since(lastUpdate) >= 1*time.Second && progressPercentage != lastPercentage {
							progress := createProgressBar(progressPercentage)
							e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
								SetContent(fmt.Sprintf("Converting video\n%s %.2f%%", progress, progressPercentage)).
								Build())
							lastUpdate = time.Now()
							lastPercentage = progressPercentage
						}
					}
				}
			}
		}

		if operation == "curldownload" {
			line := scanner.Text()
			if strings.HasPrefix(line, "###") {
				percentage := math.Min(float64(strings.Count(line, "#"))*2, 100.0)

				if time.Since(lastUpdate) >= updateInterval && percentage != lastPercentage {
					progress := createProgressBar(percentage)
					e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
						SetContent(fmt.Sprintf("Downloading video\n%s %.2f%%", progress, percentage)).
						Build())
					lastUpdate = time.Now()
					lastPercentage = percentage
				}
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.HandleError(e, "Timed out", fmt.Sprintf("%s canceled as it took too long", utils.CapitalizeFirstLetter(operation)))
			return "", fmt.Errorf("operation timed out")
		}
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("%s failed: %s", operation, errMsg)
	}

	if operation == "download" || operation == "curldownload" {
		files, err := filepath.Glob(filepath.Join(task.tempDir, "video_download.*"))
		if err != nil {
			return "", fmt.Errorf("error finding downloaded file: %w", err)
		}

		return files[0], nil
	}

	return task.filePathProcessed, nil
}

func getCodec(file string) (string, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=codec_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		file,
	)

	output, err := cmd.CombinedOutput()
	if err == nil && strings.TrimSpace(string(output)) != "" {
		return strings.TrimSpace(string(output)), nil
	}

	cmd = exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=codec_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		file,
	)

	output, err = cmd.CombinedOutput()
	if err != nil {
		slog.Error("Error running ffprobe",
			"error", err,
			"output", string(output),
			"file", file,
		)
		return "", fmt.Errorf("ffprobe error: %w, output: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}

func createProgressBar(percentage float64) string {
	filledBlocks := int(percentage / 100 * float64(20))

	bar := strings.Repeat("█", filledBlocks) + strings.Repeat("░", 20-filledBlocks)

	return bar
}

func cleanup(tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		fmt.Printf("Error while removing downloaded files: %v", err)
	}
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

func attachFile(e *handler.CommandEvent, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		utils.HandleError(e, "Error while opening file", err.Error())
		return err
	}
	defer file.Close()

	codec, err := getCodec(filePath)
	if err != nil {
		utils.HandleError(e, "Error getting codec", err.Error())
		return err
	}

	var outputFileName string
	switch codec {
	case "mp3":
		outputFileName = "audio.mp3"
	default:
		outputFileName = "video.mp4"
	}

	_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
		SetContent("").
		AddFile(outputFileName, filePath, file).
		Build())
	if err != nil {
		return err
	}

	utils.UpdateStatistics("video_downloads")

	return nil
}

func handleSpecialScenarios(e *handler.CommandEvent, url string) (string, error) {
	if strings.HasPrefix(url, "https://ylilauta.org/file/") {
		parts := strings.Split(url, "/")
		if len(parts) == 0 {
			utils.HandleError(e, "Invalid URL format", "Invalid URL format")
			return "", fmt.Errorf("invalid URL format")
		}

		fileID := parts[len(parts)-1]

		if len(fileID) < 4 {
			utils.HandleError(e, "File ID is too short", "File ID is too short")
			return "", fmt.Errorf("file ID is too short")
		}

		subPath := fmt.Sprintf("%s/%s", fileID[:2], fileID[2:4])

		newURL := fmt.Sprintf("https://i.ylilauta.org/%s/%s-apple.mp4", subPath, fileID)

		url = newURL
	}

	if strings.HasPrefix(url, "https://i.ylilauta.org/") {
		parts := strings.Split(url, "/")
		if len(parts) == 0 {
			utils.HandleError(e, "Invalid URL format", "Invalid URL format")
			return "", fmt.Errorf("invalid URL format")
		}

		filename := parts[len(parts)-1]
		if !strings.HasSuffix(filename, "-apple.mp4") {
			filename = strings.TrimSuffix(filename, ".mp4") + "-apple.mp4"

			parts[len(parts)-1] = filename
			newURL := strings.Join(parts, "/")

			url = newURL
		}
	}

	return url, nil
}
