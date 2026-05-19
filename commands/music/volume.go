package music

import (
	"context"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

var VolumeCommand = discord.SlashCommandCreate{
	Name:        "volume",
	Description: "Set the volume",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionInt{
			Name:        "volume",
			Description: "The volume to set",
			Required:    true,
		},
	},
}

func HandleVolume(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if !b.Config.Lavalink.Enabled {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("Lavalink connection has not been established"))
		}

		data := e.SlashCommandInteractionData()
		volume := data.Int("volume")

		player := b.Lavalink.ExistingPlayer(*e.GuildID())
		if player == nil {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("No player found"))
		}

		if err := player.Update(context.TODO(), lavalink.WithVolume(volume)); err != nil {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("Failed to set volume"))
		}

		return e.CreateMessage(discord.NewMessageCreate().WithContentf("Volume set to %d", volume))
	}
}
