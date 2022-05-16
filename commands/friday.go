package commands

import (
	"log"
	"math/rand"
	"time"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

var friday = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "friday",
		Description: "Post a friday celebration video",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "version",
				Description: "Which version (1-5) you want? Leave empty for random",
				Required:    false,
			},
		},
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var videos [5]string
		videos[0] = "https://cdn.discordapp.com/attachments/754470348145295360/908680252157480980/fonkymonkyfriday.mp4"
		videos[1] = "https://cdn.discordapp.com/attachments/754470348145295360/908671890111991848/fonky_monky_2.mp4"
		videos[2] = "https://cdn.discordapp.com/attachments/754470348145295360/975746878777987082/perjantai.mp4"
		videos[3] = "https://cdn.discordapp.com/attachments/754470348145295360/975746878371151882/nyt_on_perjantai.mp4"
		videos[4] = "https://cdn.discordapp.com/attachments/754470348145295360/975746876110409809/Perjantai_1.mp4"

		var video int
		if len(i.ApplicationCommandData().Options) > 0 {
			value := i.ApplicationCommandData().Options[0].IntValue()
			if value >= 1 && value <= 5 {
				video = int(i.ApplicationCommandData().Options[0].IntValue()) - 1
			} else {
				// Print out random since the value does not exist
				rand.Seed(time.Now().UnixNano())
				random := rand.Intn(len(videos))
				video = random
			}
		} else {
			rand.Seed(time.Now().UnixNano())
			random := rand.Intn(len(videos))
			video = random
		}

		if err := utils.InteractionRespondMessage(s, i, videos[video]); err != nil {
			log.Print(err)
		}
	},
}
