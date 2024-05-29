package commands

import (
	"fmt"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var queueCommand = discord.SlashCommandCreate{
	Name:        "queue",
	Description: "View the current queue",
}

func HandleQueue(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		queue := b.Queues.Get(*e.GuildID())
		if queue == nil {
			return e.CreateMessage(discord.MessageCreate{
				Content: "No player found",
			})
		}

		if len(queue.Tracks) == 0 {
			return e.CreateMessage(discord.MessageCreate{
				Content: "No tracks in queue",
			})
		}

		var tracks string
		for i, track := range queue.Tracks {
			tracks += fmt.Sprintf("%d. [`%s`](<%s>)\n", i+1, track.Info.Title, *track.Info.URI)
		}

		return e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Queue:\n%s", tracks),
		})
	}
}
