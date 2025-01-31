package commands

import (
	"context"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

var seekCommand = discord.SlashCommandCreate{
	Name:        "seek",
	Description: "Seek to a position in the current track",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionInt{
			Name:        "position",
			Description: "The position to seek to",
			Required:    true,
		},
		discord.ApplicationCommandOptionInt{
			Name:        "unit",
			Description: "The unit to seek in",
			Required:    false,
			Choices: []discord.ApplicationCommandOptionChoiceInt{
				{
					Name:  "Milliseconds",
					Value: int(lavalink.Millisecond),
				},
				{
					Name:  "Seconds",
					Value: int(lavalink.Second),
				},
				{
					Name:  "Minutes",
					Value: int(lavalink.Minute),
				},
				{
					Name:  "Hours",
					Value: int(lavalink.Hour),
				},
			},
		},
	},
}

func HandleSeek(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if !b.Config.Lavalink.Enabled {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Lavalink connection has not been established").Build())
		}

		data := e.SlashCommandInteractionData()

		player := b.Lavalink.ExistingPlayer(*e.GuildID())
		if player == nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No player found").Build())
		}

		position := data.Int("position")
		unit, ok := data.OptInt("unit")
		if !ok {
			unit = int(lavalink.Second)
		}
		finalPosition := lavalink.Duration(position * unit)
		if err := player.Update(context.TODO(), lavalink.WithPosition(finalPosition)); err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Error while seeking to position %d: %s", position, err.Error()).Build())
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Seeked to position %d", position).Build())
	}
}
