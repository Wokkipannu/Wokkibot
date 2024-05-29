package commands

import (
	"context"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

var skipCommand = discord.SlashCommandCreate{
	Name:        "skip",
	Description: "Skip the current song",
}

func HandleSkip(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		player := b.Lavalink.ExistingPlayer(*e.GuildID())
		queue := b.Queues.Get(*e.GuildID())

		if player == nil || queue == nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No player or queue found").Build())
		}

		track, ok := queue.Next()
		if !ok {
			if player != nil {
				player.Update(context.TODO(), lavalink.WithNullTrack())
				return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Skipped track, no more tracks in queue").Build())
			}

			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No tracks in queue").Build())
		}

		if err := player.Update(context.TODO(), lavalink.WithTrack(track)); err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to skip track").Build())
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Skipped track").Build())
	}
}
