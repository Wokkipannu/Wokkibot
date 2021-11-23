package commands

import (
	"fmt"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

var resume = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "resume",
		Description: "Resume current track",
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if _, found := utils.Queue[i.GuildID]; found {
			if err := Conn.SetPaused(i.GuildID, false); err != nil {
				utils.InteractionRespondMessage(s, i, fmt.Sprintf("Error when trying to resume: %v", err.Error()))
			}
			utils.InteractionRespondMessage(s, i, "Track resumed")
		} else {
			utils.InteractionRespondMessage(s, i, "Nothing to resume")
		}
	},
}
