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
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Lavalink connection has not been established").Build())
		}

		player := b.Lavalink.ExistingPlayer(*e.GuildID())
		if player == nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No player found").Build())
		}

		if err := b.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), nil, false, false); err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to disconnect").Build())
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Disconnected").Build())
	}
}
