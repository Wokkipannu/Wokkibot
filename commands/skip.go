package commands

import (
	"context"
	"fmt"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
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

		if err := Skip(b, e.GuildID()); err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(utils.CapitalizeFirstLetter(err.Error())).Build())
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Skipped track").Build())

		// track, ok := queue.Next()
		// if !ok {
		// 	player.Update(context.TODO(), lavalink.WithNullTrack())
		// 	return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Skipped track, no more tracks in queue").Build())
		// }

		// if err := player.Update(context.TODO(), lavalink.WithTrack(track)); err != nil {
		// 	return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to skip track").Build())
		// }

		// return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Skipped track").Build())
	}
}

func Skip(b *wokkibot.Wokkibot, guildId *snowflake.ID) error {
	player := b.Lavalink.ExistingPlayer(*guildId)
	queue := b.Queues.Get(*guildId)

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
