package commands

import (
	"fmt"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

var seek = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "seek",
		Description: "Seek currently playing track",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "position",
				Description: "Position to seek for",
				Required:    true,
			},
		},
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		position := uint(i.ApplicationCommandData().Options[0].UintValue())
		if position == 0 {
			utils.InteractionRespondMessage(s, i, "Value has to be higher than 0")
			return
		}

		if q, found := utils.Queue[i.GuildID]; found {
			if q.Queue[0].TrackInfo.Seekable {
				Conn.Seek(i.GuildID, position)
				utils.InteractionRespondMessage(s, i, fmt.Sprintf("Seeking from position %vs", position))
			} else {
				utils.InteractionRespondMessage(s, i, "Track is not seekable")
			}
		} else {
			utils.InteractionRespondMessage(s, i, "Nothing is playing")
		}
	},
}
