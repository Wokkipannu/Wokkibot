package commands

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type RequestPayload struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Images []string `json:"images"`
}

type ResponsePayload struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

var inspectCommand = discord.SlashCommandCreate{
	Name:        "inspect",
	Description: "Inspect an image",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "prompt",
			Description: "The prompt to use",
			Required:    true,
		},
		discord.ApplicationCommandOptionAttachment{
			Name:        "attachment",
			Description: "The image to inspect",
			Required:    true,
		},
	},
}

func HandleInspect(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		url := b.Config.AISettings.ApiUrl + "/api/generate"

		attachment := data.Attachment("attachment")
		prompt := data.String("prompt")

		base64Image, err := ImageToBase64(attachment.ProxyURL)
		if err != nil {
			log.Printf("Error converting image to base64: %v", err)
		}

		payload := RequestPayload{
			Model:  "llava",
			Prompt: prompt,
			Images: []string{base64Image},
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshaling request payload: %v", err)
			return nil
		}

		e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return nil
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("I encountered an error while trying to generate a response. Error executing request").Build())
			log.Printf("Error executing request: %v", err)
			return nil
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent("I encountered an error while trying to generate a response. Non-OK HTTP status").Build())
			log.Printf("Non-OK HTTP status: %s", resp.Status)
			return nil
		}

		var responseString string
		decoder := json.NewDecoder(resp.Body)

		var index int
		index = 0

		for {
			var responsePayload ResponsePayload
			if err := decoder.Decode(&responsePayload); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error decoding JSON:", err)
				return nil
			}

			responseString += responsePayload.Response

			index += 1

			if index%20 == 0 {
				_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent(responseString).Build())
				if err != nil {
					fmt.Println("Error updating interaction response:", err)
				}
			}

			if responsePayload.Done {
				break
			}
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetContent(responseString).Build())
		return err
	}
}

func ImageToBase64(imageURL string) (string, error) {
	response, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %v", err)
	}
	defer response.Body.Close()

	imageBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %v", err)
	}

	encodedImage := base64.StdEncoding.EncodeToString(imageBytes)

	return encodedImage, nil
}
