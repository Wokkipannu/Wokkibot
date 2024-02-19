package commands

import (
	"fmt"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
)

var resume = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "resume",
		Description: "Resume current track",
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if _, found := utils.Queue[i.GuildID]; found {
			if err := WaterlinkConnection.Guild(snowflake.MustParse(i.GuildID)).SetPaused(false); err != nil {
				utils.InteractionRespondMessage(s, i, fmt.Sprintf("Error when trying to resume: %v", err.Error()))
			}
			utils.InteractionRespondMessage(s, i, "Track resumed")
		} else {
			utils.InteractionRespondMessage(s, i, "Nothing to resume")
		}
	},
}
