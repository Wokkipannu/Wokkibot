package ping

import (
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var PingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Replies with pong!",
}

func HandlePing(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Pong!").Build())
	}
}
