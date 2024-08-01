package wokkibot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
	"wokkibot/utils"

	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"

	"github.com/mvdan/xurls"
)

type RequestPayload struct {
	Model    string    `json:"model"`
	Prompt   string    `json:"prompt"`
	Messages []Message `json:"messages"`
	System   string    `json:"system"`
}

type Message struct {
	Role    string      `json:"role"`
	Content string      `json:"content"`
	Images  interface{} `json:"images"`
}

type ResponsePayload struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   Message   `json:"message"`
	Done      bool      `json:"done"`
}

var (
	chatHistory = []Message{}
	mu          sync.Mutex
)

func (b *Wokkibot) onMessageCreate(event *events.MessageCreate) {
	HandleQuoteMessages(b, event)
	HandleCustomCommand(b, event)
	HandleXLinks(b, event)

	if b.Config.AISettings.Enabled {
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
}

func (b *Wokkibot) HandleAIResponse(e *events.MessageCreate) {
	self, _ := b.Client.Caches().SelfUser()
	input := e.Message.Content

	pattern := fmt.Sprintf(`<@%s>`, regexp.QuoteMeta(self.ID.String()))
	re := regexp.MustCompile(pattern)
	output := re.ReplaceAllString(input, "")

	mu.Lock()
	chatHistory = append(chatHistory, Message{
		Role:    "user",
		Content: e.Message.Author.EffectiveName() + " says to you \"" + output + "\"",
	})
	mu.Unlock()

	url := b.Config.AISettings.ApiUrl + "/api/chat"

	systemMessage := Message{
		Role:    "system",
		Content: b.Config.AISettings.System,
	}

	if len(chatHistory) > b.Config.AISettings.HistoryCount {
		chatHistory = chatHistory[len(chatHistory)-b.Config.AISettings.HistoryCount:]
	}

	payload := RequestPayload{
		Model:    b.Config.AISettings.Model,
		Messages: append([]Message{systemMessage}, chatHistory...),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling request payload: %v", err)
		return
	}

	msg, _ := e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().SetContent("I am thinking...").SetMessageReferenceByID(e.Message.ID).SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e.Client().Rest().UpdateMessage(e.ChannelID, msg.ID, discord.NewMessageUpdateBuilder().SetContent("I encountered an error while trying to generate a response. Error executing request").Build())
		e.Client().Rest().AddReaction(e.ChannelID, msg.ID, "❌")
		log.Printf("Error executing request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		e.Client().Rest().UpdateMessage(e.ChannelID, msg.ID, discord.NewMessageUpdateBuilder().SetContent("I encountered an error while trying to generate a response. Non-OK HTTP status").Build())
		e.Client().Rest().AddReaction(e.ChannelID, msg.ID, "❌")
		log.Printf("Non-OK HTTP status: %s", resp.Status)
		return
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
			e.Client().Rest().AddReaction(e.ChannelID, msg.ID, "❌")
			fmt.Println("Error decoding JSON:", err)
			return
		}

		responseString += responsePayload.Message.Content

		index += 1

		if index%20 == 0 {
			e.Client().Rest().UpdateMessage(e.ChannelID, msg.ID, discord.NewMessageUpdateBuilder().SetContent(responseString).Build())
		}

		if responsePayload.Done {
			break
		}
	}

	mu.Lock()
	chatHistory = append(chatHistory, Message{
		Role:    "assistant",
		Content: responseString,
	})
	mu.Unlock()

	e.Client().Rest().UpdateMessage(e.ChannelID, msg.ID, discord.NewMessageUpdateBuilder().SetContent(responseString).Build())
	e.Client().Rest().AddReaction(e.ChannelID, msg.ID, "✅")
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

func HandleXLinks(b *Wokkibot, e *events.MessageCreate) {
	self, _ := b.Client.Caches().SelfUser()
	if e.Message.Author.ID == self.ID {
		return
	}

	message := e.Message.Content

	if strings.Contains(message, "https://x.com") || strings.Contains(message, "http://x.com") {
		links := xurls.Strict.FindAllString(message, -1)

		fixedURL, err := utils.ReplaceDomain(links[0], "fixupx.com")
		if err != nil {
			return
		}

		suppressEmbeds := discord.MessageFlagSuppressEmbeds
		e.Client().Rest().UpdateMessage(e.Message.ChannelID, e.Message.ID, discord.MessageUpdate{Flags: &suppressEmbeds})

		e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().SetContent(fixedURL).SetMessageReferenceByID(e.Message.ID).SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())
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
			e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().SetContent(cmd.Output).SetMessageReferenceByID(e.Message.ID).SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())
		}
	}
}
