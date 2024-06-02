package commands

import (
	"fmt"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

var queueCommand = discord.SlashCommandCreate{
	Name:        "queue",
	Description: "View the current queue",
}

func HandleQueue(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		queue := b.Queues.Get(*e.GuildID())
		player := b.Lavalink.ExistingPlayer(*e.GuildID())
		if queue == nil || player == nil {
			return e.CreateMessage(discord.MessageCreate{
				Content: "No player found",
			})
		}

		embed := createResponseEmbed(queue, player)
		if len(queue.Tracks) > 0 || player.Track() != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).AddActionRow(discord.NewPrimaryButton("Skip", "/queue/skip").WithEmoji(discord.ComponentEmoji{Name: "⏩"})).Build())
		}
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
	}
}

func HandleQueueSkipAction(b *wokkibot.Wokkibot, e *handler.ComponentEvent) error {
	queue := b.Queues.Get(*e.GuildID())
	player := b.Lavalink.ExistingPlayer(*e.GuildID())
	if queue == nil || player == nil {
		return e.Respond(discord.InteractionResponseTypeUpdateMessage, discord.NewMessageUpdateBuilder().SetContent("No player found").Build())
	}

	err := Skip(b, e.GuildID())

	embed := createResponseEmbed(queue, player)

	var content string
	if err != nil {
		content = utils.CapitalizeFirstLetter(err.Error())
	} else {
		content = "Skipped track"
	}

	if len(queue.Tracks) > 0 || player.Track() != nil {
		return e.Respond(discord.InteractionResponseTypeUpdateMessage, discord.NewMessageUpdateBuilder().SetContent(content).SetEmbeds(embed.Build()).AddActionRow(discord.NewPrimaryButton("Skip", "/queue/skip").WithEmoji(discord.ComponentEmoji{Name: "⏩"})).Build())
	} else {
		return e.Respond(discord.InteractionResponseTypeUpdateMessage, discord.NewMessageUpdateBuilder().SetContent(content).SetEmbeds(embed.Build()).ClearContainerComponents().Build())
	}
}

func createResponseEmbed(queue *wokkibot.Queue, player disgolink.Player) *discord.EmbedBuilder {
	embed := discord.NewEmbedBuilder().SetTitle("Queue")
	embed.SetColor(utils.RGBToInteger(255, 215, 0))
	currentTrack := player.Track()

	if currentTrack != nil {
		embed.AddField("Current track", fmt.Sprintf("[%s](<%s>)", currentTrack.Info.Title, *currentTrack.Info.URI), true)
		embed.AddField("Source", currentTrack.Info.SourceName, true)
		embed.AddField("Position", fmt.Sprintf("%s / %s", utils.FormatPosition(player.Position()), utils.FormatPosition(currentTrack.Info.Length)), true)
	}

	if len(queue.Tracks) == 0 {
		if currentTrack == nil {
			embed.SetDescription("No tracks in queue")
		}
	} else {
		var tracks string
		var sources string

		for i, track := range queue.Tracks {
			tracks += fmt.Sprintf("%v. [%s](<%s>)\n", i+1, track.Info.Title, *track.Info.URI)
			sources += track.Info.SourceName + "\n"
		}

		embed.AddField("Track", tracks, true)
		embed.AddField("Source", sources, true)
	}

	return embed
}
