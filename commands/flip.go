package commands

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

var flip = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "flip",
		Description: "Flip a coin",
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		rand.Seed(time.Now().UnixNano())
		value := rand.Intn(2)
		var side string
		if value == 1 {
			side = "Heads"
		} else {
			side = "Tails"
		}
		if err := utils.InteractionRespondMessage(s, i, fmt.Sprintf("Flipped **%v**", side)); err != nil {
			log.Print(err)
		}
	},
}
