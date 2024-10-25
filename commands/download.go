package commands

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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
	filePath          string
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
		filename := fmt.Sprintf("%v.mp4", randomFileName)
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

		filePath := filepath.Join(tempDir, filename)
		filePathProcessed := filepath.Join(tempDir, processedFilename)

		task := DownloadTask{
			e:                 e,
			url:               data.String("url"),
			filePath:          filePath,
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

	e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Downloading video...").Build())
	cmd := exec.Command("yt-dlp", task.url, "-o", task.filePath)
	if err := cmd.Run(); err != nil {
		e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while downloading video").Build())
		return
	}

	e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Converting video...").Build())
	conversion := exec.Command("ffmpeg",
		"-i", task.filePath,
		"-c:v", "libx265",
		"-c:a", "aac",
		"-pix_fmt", "yuv420p",
		"-f", "mp4",
		task.filePathProcessed,
	)

	output, err := conversion.CombinedOutput()
	if err != nil {
		e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while processing video").Build())
		return
	} else {
		fmt.Println(string(output))
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

	if err := os.RemoveAll(task.tempDir); err != nil {
		fmt.Printf("Error while removing downloaded files: %v", err)
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
