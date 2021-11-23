package commands

import (
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

			for _, track := range q.Queue {
				names = append(names, track.Requester.Nick)
				tracks = append(tracks, track.TrackInfo.Title)
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

			utils.InteractionRespondMessageEmbed(s, i, embed)
		} else {
			utils.InteractionRespondMessage(s, i, "No queue found")
		}
	},
}
