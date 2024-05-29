package commands

import (
	"context"
	"fmt"
	"regexp"
	"time"
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

var playCommand = discord.SlashCommandCreate{
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
				{
					Name:  "Deezer",
					Value: "dzsearch",
				},
				{
					Name:  "Deezer ISRC",
					Value: "dzisrc",
				},
				{
					Name:  "Spotify",
					Value: "spsearch",
				},
				{
					Name:  "AppleMusic",
					Value: "amsearch",
				},
			},
		},
	},
}

func HandlePlay(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		queue := b.Queues.Get(*e.GuildID())

		identifier := data.String("identifier")

		if source, ok := data.OptString("source"); ok {
			identifier = lavalink.SearchType(source).Apply(identifier)
		} else if !urlPattern.MatchString(identifier) && !searchPattern.MatchString(identifier) {
			identifier = lavalink.SearchTypeYouTube.Apply(identifier)
		}

		voiceState, ok := b.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)
		if !ok {
			return e.CreateMessage(discord.MessageCreate{
				Content: "You need to be in a voice channel to use this command",
			})
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
				e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Nothing found for: `%s`", identifier).Build())
			},
			func(err error) {
				e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Error while looking up query: `%s`", identifier).Build())
			},
		))

		if toPlay == nil {
			return nil
		}

		if err := b.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), voiceState.ChannelID, false, false); err != nil {
			return err
		}

		if len(queue.Tracks) == 0 && b.Lavalink.ExistingPlayer(*e.GuildID()) == nil {
			e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitle("Playing").SetDescription(fmt.Sprintf("Playing: [`%s`](<%s>)", toPlay.Info.Title, *toPlay.Info.URI)).Build()).Build())

			return b.Lavalink.Player(*e.GuildID()).Update(context.TODO(), lavalink.WithTrack(*toPlay))
		}

		e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitle("Queued").SetDescription(fmt.Sprintf("Queued: [`%s`](<%s>)", toPlay.Info.Title, *toPlay.Info.URI)).Build()).Build())
		queue.Add(*toPlay)
		return nil
	}
}
