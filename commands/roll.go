package commands

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

var roll = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "roll",
		Description: "Roll a die",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "max",
				Description: "Highest possible roll value",
				Required:    false,
			},
		},
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		rand.Seed(time.Now().UnixNano())
		min := int64(1)
		var max int64
		if len(i.ApplicationCommandData().Options) > 0 {
			max = i.ApplicationCommandData().Options[0].IntValue()
			if max < 2 {
				if err := utils.InteractionRespondMessage(s, i, "Value has to be higher than 1"); err != nil {
					log.Print(err)
				}
				return
			}
		} else {
			max = int64(100)
		}
		value := min + rand.Int63n(max-min)
		if err := utils.InteractionRespondMessage(s, i, fmt.Sprintf("%v (1-%v)", value, max)); err != nil {
			log.Print(err)
		}
	},
}
