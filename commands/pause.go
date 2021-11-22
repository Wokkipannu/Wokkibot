package commands

import (
	"fmt"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

var pause = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "pause",
		Description: "Pause current track",
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if q, found := utils.Queue[i.GuildID]; found {
			if len(q.Queue) > 0 {
				if err := Conn.SetPaused(i.GuildID, true); err != nil {
					utils.InteractionRespondMessage(s, i, fmt.Sprintf("Error when trying to pause: %v", err.Error()))
				}
				utils.InteractionRespondMessage(s, i, "Track paused")
			} else {
				utils.InteractionRespondMessage(s, i, "Nothing to pause")
			}
		} else {
			utils.InteractionRespondMessage(s, i, "Nothing to pause")
		}
	},
}
