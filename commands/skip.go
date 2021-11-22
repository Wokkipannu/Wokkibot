package commands

import (
	"fmt"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

var skip = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "skip",
		Description: "Skip current track",
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if q, found := utils.Queue[i.GuildID]; found {
			Conn.Stop(i.GuildID)
			utils.InteractionRespondMessage(s, i, fmt.Sprintf("\"%v\" skipped", q.Queue[0].TrackInfo.Title))
		} else {
			utils.InteractionRespondMessage(s, i, "Nothing to skip")
		}
	},
}
