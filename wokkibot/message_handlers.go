package wokkibot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"wokkibot/utils"

	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"

	"github.com/mvdan/xurls"
)

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var (
	chatHistory = []ChatMessage{}
	mu          sync.Mutex
)

func (b *Wokkibot) onMessageCreate(event *events.MessageCreate) {
	HandleQuoteMessages(b, event)
	HandleCustomCommand(b, event)

	self, _ := b.Client.Caches().SelfUser()
	if event.Message.Author.ID == self.ID {
		return
	}
	for _, user := range event.Message.Mentions {
		if user.ID == self.ID {
			HandleAIResponse(b, event)
		}
	}
}

func HandleAIResponse(b *Wokkibot, e *events.MessageCreate) {
	if len(chatHistory) == 0 {
		chatHistory = append(chatHistory, ChatMessage{
			Role:    "system",
			Content: b.Config.OpenAIInstructions,
		})
	}

	msg := strings.TrimPrefix(e.Message.Content, "<@512004300218695714> ")

	userMessage := e.Message.Member.EffectiveName() + " (" + e.Message.Member.User.ID.String() + ") said to you: " + msg

	mu.Lock()
	chatHistory = append(chatHistory, ChatMessage{
		Role:    "user",
		Content: userMessage,
	})
	mu.Unlock()

	response, err := getOpenAIResponse(b.Config.OpenAIApiKey)
	if err != nil {
		log.Printf("Error getting OpenAI response: %v", err)
		e.Client().Rest().CreateMessage(e.ChannelID, discord.NewMessageCreateBuilder().SetContent("Sorry, I encountered an error.").Build())
		return
	}

	mu.Lock()
	chatHistory = append(chatHistory, ChatMessage{
		Role:    "assistant",
		Content: response,
	})
	mu.Unlock()

	e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().SetContent(response).SetMessageReferenceByID(e.Message.ID).Build())
}

func getOpenAIResponse(apiKey string) (string, error) {
	mu.Lock()
	history := chatHistory
	mu.Unlock()

	openAIPrompt := map[string]interface{}{
		"model":    "gpt-4-turbo",
		"messages": history,
	}

	requestBody, err := json.Marshal(openAIPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK HTTP status: %s, body: %s", resp.Status, string(body))
	}

	var openAIResp OpenAIResponse
	err = json.Unmarshal(body, &openAIResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	if len(openAIResp.Choices) > 0 {
		return openAIResp.Choices[0].Message.Content, nil
	}

	return "", nil
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
