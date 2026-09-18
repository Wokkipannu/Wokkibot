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

var (
	xLinkPattern = regexp.MustCompile(`https?:\/\/(x|twitter)\.com\/(.*\/status\/\d+)\??.*`)
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

/*
 * Since embeds might appear after we've already suppressed them, we need to check for them again on update
 */
func (h *Handler) OnMessageUpdate(e *events.MessageUpdate) {
	if e.GuildID != nil {
		h.EnsureGuildExists(*e.GuildID)
	}

	if xLinkPattern.MatchString(e.Message.Content) {
		if guild, exists := h.Guilds[*e.GuildID]; exists && guild.ConvertXLinks {
			// Prevent update loop, only suppress if not already suppressed
			if !e.Message.Flags.Has(discord.MessageFlagSuppressEmbeds) {
				h.SuppressEmbeds(nil, e)
			}
		}
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
		msg, err := e.Client().Rest.GetMessage(channelId, messageId)
		if err != nil {
			return
		}

		embed := utils.QuoteEmbed(*msg)

		e.Client().Rest.CreateMessage(e.Message.ChannelID, discord.NewMessageCreate().WithEmbeds(embed).AddActionRow(discord.NewLinkButton("Go to message", links[0])))
	}
}

func (h *Handler) HandleXLinks(e *events.MessageCreate) {
	self, _ := e.Client().Caches.SelfUser()
	if e.Message.Author.ID == self.ID {
		return
	}

	message := e.Message.Content

	if xLinkPattern.MatchString(message) {
		links := xurls.Strict.FindAllString(message, -1)

		fixedURL, err := utils.ReplaceDomain(links[0], "fixvx.com")
		if err != nil {
			return
		}

		h.SuppressEmbeds(e, nil)

		e.Client().Rest.CreateMessage(e.Message.ChannelID, discord.NewMessageCreate().WithContent(fixedURL).WithMessageReferenceByID(e.Message.ID).WithAllowedMentions(&discord.AllowedMentions{RepliedUser: false}))
	}
}

func (h *Handler) SuppressEmbeds(eC *events.MessageCreate, eU *events.MessageUpdate) {
	suppressEmbeds := discord.MessageFlagSuppressEmbeds

	if eC != nil {
		eC.Client().Rest.UpdateMessage(eC.Message.ChannelID, eC.Message.ID, discord.MessageUpdate{
			Flags: &suppressEmbeds,
		})
	}
	if eU != nil {
		eU.Client().Rest.UpdateMessage(eU.Message.ChannelID, eU.Message.ID, discord.MessageUpdate{
			Flags: &suppressEmbeds,
		})
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

			e.Client().Rest.CreateMessage(e.Message.ChannelID, discord.NewMessageCreate().
				WithContent(output).
				WithMessageReferenceByID(e.Message.ID).
				WithAllowedMentions(&discord.AllowedMentions{RepliedUser: false}))
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
