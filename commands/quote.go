package commands

import (
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var quoteCommand = discord.MessageCommandCreate{
	Name: "Quote",
}

func HandleQuote(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		msg := e.MessageCommandInteractionData().TargetMessage()

		embed := utils.QuoteEmbed(msg)

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).AddActionRow(discord.NewLinkButton("Go to message", msg.JumpURL())).Build())
	}
}
