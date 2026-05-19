package music

import (
	"fmt"
	"wokkibot/queue"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

var QueueCommand = discord.SlashCommandCreate{
	Name:        "queue",
	Description: "View the current queue",
}

func HandleQueue(b *wokkibot.Wokkibot, q *queue.QueueManager) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if !b.Config.Lavalink.Enabled {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("Lavalink connection has not been established"))
		}

		queue := q.Get(*e.GuildID())
		player := b.Lavalink.ExistingPlayer(*e.GuildID())
		if queue == nil || player == nil {
			return e.CreateMessage(discord.MessageCreate{
				Content: "No player found",
			})
		}

		embed := createResponseEmbed(queue, player)
		if len(queue.Tracks) > 0 || player.Track() != nil {
			return e.CreateMessage(discord.NewMessageCreate().WithEmbeds(embed).AddActionRow(discord.NewPrimaryButton("Skip", "/queue/skip").WithEmoji(discord.ComponentEmoji{Name: "⏩"})))
		}
		return e.CreateMessage(discord.NewMessageCreate().WithEmbeds(embed))
	}
}

func HandleQueueSkipAction(b *wokkibot.Wokkibot, q *queue.QueueManager, e *handler.ComponentEvent) error {
	queue := q.Get(*e.GuildID())
	player := b.Lavalink.ExistingPlayer(*e.GuildID())
	if queue == nil || player == nil {
		return e.Respond(discord.InteractionResponseTypeUpdateMessage, discord.NewMessageUpdate().WithContent("No player found"))
	}

	err := Skip(b, q, e.GuildID())

	embed := createResponseEmbed(queue, player)

	var content string
	if err != nil {
		content = utils.CapitalizeFirstLetter(err.Error())
	} else {
		content = "Skipped track"
	}

	if len(queue.Tracks) > 0 || player.Track() != nil {
		return e.Respond(discord.InteractionResponseTypeUpdateMessage, discord.NewMessageUpdate().WithContent(content).WithEmbeds(embed).AddActionRow(discord.NewPrimaryButton("Skip", "/queue/skip").WithEmoji(discord.ComponentEmoji{Name: "⏩"})))
	} else {
		return e.Respond(discord.InteractionResponseTypeUpdateMessage, discord.NewMessageUpdate().WithContent(content).WithEmbeds(embed).ClearComponents())
	}
}

func createResponseEmbed(queue *queue.Queue, player disgolink.Player) discord.Embed {
	embed := discord.NewEmbed().WithTitle("Queue")
	embed = embed.WithColor(utils.RGBToInteger(255, 215, 0))
	currentTrack := player.Track()

	if currentTrack != nil {
		embed = embed.AddField("Current track", fmt.Sprintf("[%s](<%s>)", currentTrack.Info.Title, *currentTrack.Info.URI), true)
		embed = embed.AddField("Source", currentTrack.Info.SourceName, true)
		embed = embed.AddField("Position", fmt.Sprintf("%s / %s", utils.FormatPosition(player.Position()), utils.FormatPosition(currentTrack.Info.Length)), true)
	}

	if len(queue.Tracks) == 0 {
		if currentTrack == nil {
			embed = embed.WithDescription("No tracks in queue")
		}
	} else {
		var tracks string
		var sources string

		for i, track := range queue.Tracks {
			tracks += fmt.Sprintf("%v. [%s](<%s>)\n", i+1, track.Info.Title, *track.Info.URI)
			sources += track.Info.SourceName + "\n"
		}

		embed = embed.AddField("Track", tracks, true)
		embed = embed.AddField("Source", sources, true)
	}

	return embed
}

func HandleQueueSkipActionComponent(b *wokkibot.Wokkibot, q *queue.QueueManager) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		HandleQueueSkipAction(b, q, e)
		return nil
	}
}
