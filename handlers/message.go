package handlers

import (
	"log/slog"
	"regexp"
	"strings"
	"time"
	"wokkibot/utils"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/mvdan/xurls"
	"golang.org/x/exp/rand"
)

func (h *Handler) OnMessageCreate(e *events.MessageCreate) {
	h.HandleQuoteMessages(e)
	h.HandleCustomCommand(e)
	h.HandleXLinks(e)
}

func (h *Handler) HandleQuoteMessages(e *events.MessageCreate) {
	prefix := "https://discord.com/channels/"
	message := e.Message.Content

	if strings.Contains(message, prefix) {
		links := xurls.Strict.FindAllString(message, -1)

		slashes := strings.Split(links[0], "/")

		channelId := snowflake.MustParse(slashes[len(slashes)-2])
		messageId := snowflake.MustParse(slashes[len(slashes)-1])
		msg, err := e.Client().Rest().GetMessage(channelId, messageId)
		if err != nil {
			return
		}

		embed := utils.QuoteEmbed(*msg)

		e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).AddActionRow(discord.NewLinkButton("Go to message", links[0])).Build())
	}
}

func (h *Handler) HandleXLinks(e *events.MessageCreate) {
	self, _ := e.Client().Caches().SelfUser()
	if e.Message.Author.ID == self.ID {
		return
	}

	message := e.Message.Content

	regexPattern := `https?:\/\/(x|twitter)\.com\/(.*\/status\/\d+)\??.*`
	r := regexp.MustCompile(regexPattern)

	if r.MatchString(message) {
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

func (h *Handler) HandleCustomCommand(e *events.MessageCreate) {
	input := e.Message.Content

	if input == "" {
		return
	}
	prefix := string(input[0])
	name := strings.TrimPrefix(input, prefix)

	for _, cmd := range h.CustomCommands {
		if cmd.Prefix == prefix && cmd.Name == name && cmd.GuildID == *e.GuildID {
			output := handleVariables(cmd.Output)

			e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().
				SetContent(output).
				SetMessageReferenceByID(e.Message.ID).
				SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).
				Build())
		}
	}
}

func handleVariables(text string) string {
	re := regexp.MustCompile(`\{\{(\w+)\|([^}]+)\}\}`)

	return re.ReplaceAllStringFunc(text, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		variable := parts[1]
		value := parts[2]

		switch variable {
		case "time":
			loc, err := time.LoadLocation(value)
			if err != nil {
				slog.Error("Failed to load timezone", "location", value, "error", err)
				return "INVALID TIMEZONE NAME"
			}
			return time.Now().In(loc).Format("15:04 MST")
		case "random":
			choices := strings.Split(value, ";")
			if len(choices) == 0 {
				return "NO CHOICES PROVIDED"
			}
			randomIndex := rand.Intn(len(choices))
			return strings.TrimSpace(choices[randomIndex])
		default:
			return match
		}
	})
}
