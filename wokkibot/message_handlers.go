package wokkibot

import (
	"regexp"
	"wokkibot/utils"

	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"

	"github.com/mvdan/xurls"
)

func (b *Wokkibot) onMessageCreate(event *events.MessageCreate) {
	HandleQuoteMessages(b, event)
	HandleCustomCommand(b, event)
	HandleXLinks(b, event)
}

func HandleQuoteMessages(b *Wokkibot, e *events.MessageCreate) {
	prefix := "https://discord.com/channels/"
	message := e.Message.Content

	if strings.Contains(message, prefix) {
		links := xurls.Strict.FindAllString(message, -1)

		slashes := strings.Split(links[0], "/")

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
