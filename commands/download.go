package commands

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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

func HandleDownload(b *wokkibot.Wokkibot) handler.CommandHandler {
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

		tempDir, err := os.MkdirTemp("downloads", "video_*")
		if err != nil {
			fmt.Printf("Failed to create temp directory: %v\n", err)
			return err
		}

		filePath := filepath.Join(tempDir, filename)
		filePathProcessed := filepath.Join(tempDir, processedFilename)

		cmd := exec.Command("yt-dlp",
			data.String("url"),
			"-o", filePath,
		)
		if err := cmd.Run(); err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while downloading video").Build())
		}

		conversion := exec.Command("ffmpeg",
			"-i", filePath,
			"-c:v", "libx264",
			"-c:a", "aac",
			"-pix_fmt", "yuv420p",
			"-f", "mp4",
			filePathProcessed,
		)

		output, cErr := conversion.CombinedOutput()
		if cErr != nil {
			_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while processing video").Build())
			return err
		} else {
			fmt.Println(string(output))
		}

		file, err := os.Open(filePathProcessed)
		if err != nil {
			e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while opening file").Build())
		}
		defer file.Close()

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().AddFile(processedFilename, processedFilename, file).Build())
		if err != nil {
			e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("Error while attaching file").Build())
		}

		if err := os.RemoveAll(tempDir); err != nil {
			fmt.Printf("Error while removing downloaded files: %v", err)
		}

		return err
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
