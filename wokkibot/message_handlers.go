package wokkibot

import (
	"wokkibot/utils"

	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"

	"github.com/mvdan/xurls"
)

func (b *Wokkibot) onMessageCreate(event *events.MessageCreate) {
	prefix := "https://discord.com/channels/"

	message := event.Message.Content

	if strings.Contains(message, prefix) {
		links := xurls.Strict.FindAllString(message, -1)

		slashes := strings.Split(links[0], "/")

		messageId := slashes[len(slashes)-1]
		msg, err := b.Client.Rest().GetMessage(event.Message.ChannelID, snowflake.MustParse(messageId))
		if err != nil {
			return
		}

		embed := utils.QuoteEmbed(*msg)

		event.Client().Rest().CreateMessage(event.Message.ChannelID, discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).AddActionRow(discord.NewLinkButton("Go to message", links[0])).Build())
	}
}
