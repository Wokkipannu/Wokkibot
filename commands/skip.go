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
			if q.Queue[0].Interaction != nil {
				Session.InteractionResponseEdit(Session.State.User.ID, q.Queue[0].Interaction, &discordgo.WebhookEdit{
					Content:    "",
					Embeds:     []*discordgo.MessageEmbed{q.Queue[0].Embed},
					Components: []discordgo.MessageComponent{},
				})
			}
			if q.Queue[0].Message != nil {
				content := ""
				Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
					ID:         q.Queue[0].Message.ID,
					Channel:    q.Queue[0].Message.ChannelID,
					Content:    &content, // Why the fuck is this *string?
					Embeds:     []*discordgo.MessageEmbed{q.Queue[0].Embed},
					Components: []discordgo.MessageComponent{},
				})
			}
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
