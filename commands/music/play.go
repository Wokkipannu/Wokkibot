package music

import (
	"context"
	"fmt"
	"regexp"
	"time"
	"wokkibot/queue"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

var (
	urlPattern    = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)
)

var PlayCommand = discord.SlashCommandCreate{
	Name:        "play",
	Description: "Play a song",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "identifier",
			Description: "Link to the song",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "source",
			Description: "The source to search on",
			Required:    false,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "YouTube",
					Value: string(lavalink.SearchTypeYouTube),
				},
				{
					Name:  "YouTube Music",
					Value: string(lavalink.SearchTypeYouTubeMusic),
				},
				{
					Name:  "SoundCloud",
					Value: string(lavalink.SearchTypeSoundCloud),
				},
				// {
				// 	Name:  "Deezer",
				// 	Value: "dzsearch",
				// },
				// {
				// 	Name:  "Deezer ISRC",
				// 	Value: "dzisrc",
				// },
				{
					Name:  "Spotify",
					Value: "spsearch",
				},
				{
					Name:  "http",
					Value: "http",
				},
				// {
				// 	Name:  "AppleMusic",
				// 	Value: "amsearch",
				// },
			},
		},
	},
}

func HandlePlay(b *wokkibot.Wokkibot, q *queue.QueueManager) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if !b.Config.Lavalink.Enabled {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("Lavalink connection has not been established"))
		}

		data := e.SlashCommandInteractionData()

		queue := q.Get(*e.GuildID())

		identifier := data.String("identifier")

		if source, ok := data.OptString("source"); ok {
			identifier = lavalink.SearchType(source).Apply(identifier)
		} else if !urlPattern.MatchString(identifier) && !searchPattern.MatchString(identifier) {
			identifier = lavalink.SearchTypeYouTube.Apply(identifier)
		}

		voiceState, ok := b.Client.Caches.VoiceState(*e.GuildID(), e.User().ID)
		if !ok {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("You need to be in a voice channel to use this command"))
		}

		if err := e.DeferCreateMessage(false); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var toPlay *lavalink.Track
		b.Lavalink.BestNode().LoadTracksHandler(ctx, identifier, disgolink.NewResultHandler(
			func(track lavalink.Track) {
				toPlay = &track
			},
			func(playlist lavalink.Playlist) {
				toPlay = &playlist.Tracks[0]
			},
			func(tracks []lavalink.Track) {
				toPlay = &tracks[0]
			},
			func() {
				e.CreateMessage(discord.NewMessageCreate().WithContentf("Nothing found for: `%s`", identifier))
			},
			func(err error) {
				e.CreateMessage(discord.NewMessageCreate().WithContentf("Error while looking up query: `%s`", identifier))
			},
		))

		if toPlay == nil {
			e.UpdateInteractionResponse(discord.NewMessageUpdate().WithContentf("Nothing found for: `%s`", identifier))
			return nil
		}

		if err := b.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), voiceState.ChannelID, false, false); err != nil {
			e.UpdateInteractionResponse(discord.NewMessageUpdate().WithContentf("Error while updating voice state: %s", err.Error()))
			return err
		}

		embed := discord.NewEmbed()
		embed = embed.WithColor(utils.RGBToInteger(255, 215, 0))
		if toPlay.Info.ArtworkURL != nil {
			embed = embed.WithImage(*toPlay.Info.ArtworkURL)
		}
		embed = embed.WithFooterTextf("Length: %s", utils.FormatDuration(toPlay.Info.Length))
		embed = embed.AddField("Track", fmt.Sprintf("[%s](<%s>)", toPlay.Info.Title, *toPlay.Info.URI), true)
		embed = embed.AddField("Source", toPlay.Info.SourceName, true)

		player := b.Lavalink.ExistingPlayer(*e.GuildID())
		if len(queue.Tracks) == 0 && (player == nil || player.Track() == nil) {
			embed = embed.WithTitle("Playing")
			e.UpdateInteractionResponse(discord.NewMessageUpdate().WithEmbeds(embed))
			return b.Lavalink.Player(*e.GuildID()).Update(context.TODO(), lavalink.WithTrack(*toPlay))
		}

		embed = embed.WithTitlef("Queued to position %d", len(queue.Tracks)+1)
		e.UpdateInteractionResponse(discord.NewMessageUpdate().WithEmbeds((embed)))
		queue.Add(*toPlay)
		return nil
	}
}
