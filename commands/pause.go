package commands

import (
	"fmt"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
)

var pause = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "pause",
		Description: "Pause current track",
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if _, found := utils.Queue[i.GuildID]; found {
			if err := WaterlinkConnection.Guild(snowflake.MustParse(i.GuildID)).SetPaused(true); err != nil {
				utils.InteractionRespondMessage(s, i, fmt.Sprintf("Error when trying to pause: %v", err.Error()))
			}
			utils.InteractionRespondMessage(s, i, "Track paused")
		} else {
			utils.InteractionRespondMessage(s, i, "Nothing to pause")
		}
	},
}
