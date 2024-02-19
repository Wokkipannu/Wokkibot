package commands

import (
	"fmt"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
)

var volume = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "volume",
		Description: "Change player volume",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "volume",
				Description: "Volume in percentage (0-100)",
				Required:    true,
			},
		},
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		vol := uint(i.ApplicationCommandData().Options[0].UintValue())
		if q, found := utils.Queue[i.GuildID]; found {
			q.Volume = vol
			WaterlinkConnection.Guild(snowflake.MustParse(i.GuildID)).UpdateVolume(uint16(vol))
			utils.InteractionRespondMessage(s, i, fmt.Sprintf("Changed player volume for all queued songs to %v", vol))
		} else {
			utils.InteractionRespondMessage(s, i, "Queue does not exist")
		}
	},
}
