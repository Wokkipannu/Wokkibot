package music

import (
	"context"
	"fmt"
	"wokkibot/queue"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

var SkipCommand = discord.SlashCommandCreate{
	Name:        "skip",
	Description: "Skip the current song",
}

func HandleSkip(b *wokkibot.Wokkibot, q *queue.QueueManager) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if !b.Config.Lavalink.Enabled {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("Lavalink connection has not been established"))
		}

		player := b.Lavalink.ExistingPlayer(*e.GuildID())
		queue := q.Get(*e.GuildID())

		if player == nil || queue == nil {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("No player or queue found"))
		}

		if err := Skip(b, q, e.GuildID()); err != nil {
			return e.CreateMessage(discord.NewMessageCreate().WithContent(utils.CapitalizeFirstLetter(err.Error())))
		}

		return e.CreateMessage(discord.NewMessageCreate().WithContent("Skipped track"))
	}
}

func Skip(b *wokkibot.Wokkibot, q *queue.QueueManager, guildId *snowflake.ID) error {
	player := b.Lavalink.ExistingPlayer(*guildId)
	queue := q.Get(*guildId)

	if player == nil || queue == nil {
		return fmt.Errorf("no player or queue found")
	}

	track, ok := queue.Next()
	if !ok {
		player.Update(context.TODO(), lavalink.WithNullTrack())
		return fmt.Errorf("skipped track, no more tracks in queue")
	}

	if err := player.Update(context.TODO(), lavalink.WithTrack(track)); err != nil {
		return fmt.Errorf("failed to skip track")
	}

	return nil
}
