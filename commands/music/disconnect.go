package music

import (
	"context"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var DisconnectCommand = discord.SlashCommandCreate{
	Name:        "disconnect",
	Description: "Disconnect the bot from the voice channel",
}

func HandleDisconnect(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if !b.Config.Lavalink.Enabled {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("Lavalink connection has not been established"))
		}

		player := b.Lavalink.ExistingPlayer(*e.GuildID())
		if player == nil {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("No player found"))
		}

		if err := b.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), nil, false, false); err != nil {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("Failed to disconnect"))
		}

		return e.CreateMessage(discord.NewMessageCreate().WithContent("Disconnected"))
	}
}
