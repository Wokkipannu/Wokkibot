package handlers

import (
	"log/slog"
	"regexp"
	"strings"
	"time"
	"wokkibot/utils"

	"math/rand/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/mvdan/xurls"
)

func (h *Handler) OnMessageCreate(e *events.MessageCreate) {
	if e.GuildID != nil {
		h.EnsureGuildExists(*e.GuildID)
	}

	h.HandleQuoteMessages(e)
	h.HandleCustomCommand(e)

	if guild, exists := h.Guilds[*e.GuildID]; exists && guild.ConvertXLinks {
		h.HandleXLinks(e)
	}
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

		fixedURL, err := utils.ReplaceDomain(links[0], "fixvx.com")
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
			output := handleVariables(cmd.Output, e)

			e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().
				SetContent(output).
				SetMessageReferenceByID(e.Message.ID).
				SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).
				Build())
		}
	}
}

func handleVariables(text string, e *events.MessageCreate) string {
	re := regexp.MustCompile(`\{\{(\w+)\|([^}]+)\}\}`)

	author := e.Message.Author

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
			randomIndex := rand.IntN(len(choices))
			return strings.TrimSpace(choices[randomIndex])
		case "user":
			switch value {
			case "name":
				return author.Username
			case "id":
				return author.ID.String()
			case "avatar":
				return author.EffectiveAvatarURL()
			case "mention":
				return author.Mention()
			case "created":
				return author.ID.Time().Format("2006-01-02 15:04:05")
			default:
				return "INVALID USER ATTRIBUTE"
			}
		default:
			return match
		}
	})
}
