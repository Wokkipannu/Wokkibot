package commands

import (
	"fmt"
	"log"
	"strings"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

var queue = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "queue",
		Description: "Display queued tracks",
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if q, found := utils.Queue[i.GuildID]; found {
			color := Session.State.UserColor(Session.State.User.ID, q.TextChannelID)

			embed := &discordgo.MessageEmbed{}
			embed.Color = color
			embed.Title = "Queue"

			var names []string
			var tracks []string
			var durations []string

			for _, track := range q.Queue {
				duration := track.TrackInfo.Length
				seconds := (duration / 1000) % 60
				minutes := (duration / (1000 * 60) % 60)
				hours := (duration / (1000 * 60 * 60) % 24)

				names = append(names, utils.EscapeString(utils.GetName(track.Requester)))
				tracks = append(tracks, fmt.Sprintf("[%v](%v)", utils.TruncateString(utils.EscapeString(track.TrackInfo.Title), 50), track.TrackInfo.URI))
				if track.TrackInfo.Stream {
					durations = append(durations, "Stream")
				} else {
					durations = append(durations, fmt.Sprintf("%v:%v:%v", utils.NumberFormat(hours), utils.NumberFormat(minutes), utils.NumberFormat(seconds)))
				}
			}

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Title",
				Value:  strings.Join(tracks, "\n"),
				Inline: true,
			})
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Requester",
				Value:  strings.Join(names, "\n"),
				Inline: true,
			})
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Duration",
				Value:  strings.Join(durations, "\n"),
				Inline: true,
			})

			if err := utils.InteractionRespondMessageEmbed(s, i, embed); err != nil {
				log.Print(err)
			}
		} else {
			if err := utils.InteractionRespondMessage(s, i, "No queue found"); err != nil {
				log.Print(err)
			}
		}
	},
}
