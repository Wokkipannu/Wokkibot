package commands

import (
	"fmt"
	"log"
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
			if err := utils.InteractionRespondMessage(s, i, fmt.Sprintf("Track \"%v\" skipped", utils.EscapeString(q.Queue[0].TrackInfo.Title))); err != nil {
				log.Print(err)
			}
		} else {
			if err := utils.InteractionRespondMessage(s, i, "Nothing to skip"); err != nil {
				log.Print(err)
			}
		}
	},
}
