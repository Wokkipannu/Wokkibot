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

		embed := discord.NewEmbedBuilder().SetTitle("Queue")
		currentTrack := b.Lavalink.ExistingPlayer(*e.GuildID()).Track()

		if len(queue.Tracks) == 0 {
			embed.AddField("", "No tracks in queue", false)
		} else {
			for i, track := range queue.Tracks {
				embed.AddField("", fmt.Sprintf("%v. [%s](<%s>)", i+1, track.Info.Title, *track.Info.URI), false)
			}
		}

		embed.SetFooterTextf("Currently playing: %s (%s)", currentTrack.Info.Title, *currentTrack.Info.URI)
		if currentTrack.Info.ArtworkURL != nil {
			embed.SetFooterIcon(*currentTrack.Info.ArtworkURL)
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
	}
}
