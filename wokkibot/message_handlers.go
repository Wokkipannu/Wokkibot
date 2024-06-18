package wokkibot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"wokkibot/utils"

	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"

	"github.com/mvdan/xurls"
)

type RequestPayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ResponsePayload struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

func (b *Wokkibot) onMessageCreate(event *events.MessageCreate) {
	HandleQuoteMessages(b, event)
	HandleCustomCommand(b, event)

	self, _ := b.Client.Caches().SelfUser()
	if event.Message.Author.ID == self.ID {
		return
	}
	for _, user := range event.Message.Mentions {
		if user.ID == self.ID {
			b.HandleAIResponse(event)
		}
	}
}

func (b *Wokkibot) HandleAIResponse(e *events.MessageCreate) {
	url := b.Config.AIApiUrl
	payload := RequestPayload{
		Model:  "llama3",
		Prompt: e.Message.Content,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling request payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-OK HTTP status: %s", resp.Status)
		return
	}

	var responseString string
	decoder := json.NewDecoder(resp.Body)

	var index int
	index = 0

	msg, _ := e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().SetContent("...").SetMessageReferenceByID(e.Message.ID).Build())

	for {
		var responsePayload ResponsePayload
		if err := decoder.Decode(&responsePayload); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		responseString += responsePayload.Response

		index += 1

		if index%5 == 0 {
			e.Client().Rest().UpdateMessage(e.ChannelID, msg.ID, discord.NewMessageUpdateBuilder().SetContent(responseString).Build())
		}

		if responsePayload.Done {
			break
		}
	}

	e.Client().Rest().UpdateMessage(e.ChannelID, msg.ID, discord.NewMessageUpdateBuilder().SetContent(responseString).Build())
}

func HandleQuoteMessages(b *Wokkibot, e *events.MessageCreate) {
	prefix := "https://discord.com/channels/"
	message := e.Message.Content

	if strings.Contains(message, prefix) {
		links := xurls.Strict.FindAllString(message, -1)

		slashes := strings.Split(links[0], "/")

		// guildId := snowflake.MustParse(slashes[len(slashes)-3])

		// if guildId != *event.Message.GuildID {
		// 	return
		// }

		channelId := snowflake.MustParse(slashes[len(slashes)-2])
		messageId := snowflake.MustParse(slashes[len(slashes)-1])
		msg, err := b.Client.Rest().GetMessage(channelId, messageId)
		if err != nil {
			return
		}

		embed := utils.QuoteEmbed(*msg)

		e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).AddActionRow(discord.NewLinkButton("Go to message", links[0])).Build())
	}
}

func HandleCustomCommand(b *Wokkibot, e *events.MessageCreate) {
	input := e.Message.Content

	if input == "" {
		return
	}
	prefix := string(input[0])
	name := strings.TrimPrefix(input, prefix)

	for _, cmd := range b.CustomCommands {
		if cmd.Prefix == prefix && cmd.Name == name && cmd.GuildID == *e.GuildID {
			e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().SetContent(cmd.Output).SetMessageReferenceByID(e.Message.ID).Build())
		}
	}
}
