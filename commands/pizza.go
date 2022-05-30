package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
	"wokkibot/config"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

type ToppingsResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

var pizza = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "pizza",
		Description: "Get random pizza toppings",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Amount of random toppings",
				Required:    false,
			},
		},
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var amount int
		if len(i.ApplicationCommandData().Options) > 0 {
			amount = int(i.ApplicationCommandData().Options[0].IntValue())
		} else {
			amount = 4
		}

		rand.Seed(time.Now().UnixNano())
		API := config.Config("PIZZAAPI")
		res, err := s.Client.Get(API + "?amount=" + strconv.Itoa(amount))
		if err != nil {
			log.Print(err)
			utils.InteractionRespondMessage(s, i, "Something went wrong when fetching toppings")
			return
		}
		defer res.Body.Close()

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Print(readErr)
			utils.InteractionRespondMessage(s, i, "Something went wrong when reading response data")
			return
		}

		toppings := ToppingsResponse{}
		jsonErr := json.Unmarshal(body, &toppings)
		if jsonErr != nil {
			log.Print(jsonErr)
			utils.InteractionRespondMessage(s, i, "Something went wrong when parsing response data")
			return
		}

		var bases [3]string
		bases[0] = "Normaali"
		bases[1] = "Runsaskuituinen"
		bases[2] = "Gluteeniton"
		base := bases[rand.Intn(len(bases))]

		var sauces [12]string
		sauces[0] = "Tomaattikastike"
		sauces[1] = "Mexicana-kastike"
		sauces[2] = "Cheddarjuustokastike"
		sauces[3] = "Barbecuekastike"
		sauces[4] = "Kebabkastike"
		sauces[5] = "Koskenlaskija-juustokastike"
		sauces[6] = "Piparjuurimajoneesi"
		sauces[7] = "Chipotlemajoneesi"
		sauces[8] = "Hotti-kastike"
		sauces[9] = "Vegaaninen valkosipulimajoneesi"
		sauces[10] = "Chilimajoneesi"
		sauces[11] = "Valkosipulimajoneesi"
		sauce := sauces[rand.Intn(len(sauces))]

		embed := &discordgo.MessageEmbed{}
		embed.Color = s.State.UserColor(s.State.User.ID, i.ChannelID)
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%v Pizza", utils.GetName(i.Member)),
			IconURL: i.Member.User.AvatarURL(""),
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Pohja",
			Value:  base,
			Inline: true,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Kastike",
			Value:  sauce,
			Inline: true,
		})

		// If the length of toppings is higher than 900, we split it to multiple embeds
		runeLenth := []rune(toppings.Data)
		var outputs []string
		if len(runeLenth) > 900 {
			outputs = append(outputs, string(runeLenth[:900]))
			outputs = append(outputs, string(runeLenth[900:]))
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Täytteet (1/2)",
				Value:  outputs[0],
				Inline: true,
			})
			utils.InteractionRespondMessageEmbed(s, i, embed)

			// Remove last field of the embed and create a new field with rest of the toppings
			embed.Fields = embed.Fields[:len(embed.Fields)-3]
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Täytteet (2/2)",
				Value:  outputs[1],
				Inline: true,
			})
			s.ChannelMessageSendEmbed(i.ChannelID, embed)
		} else {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Täytteet",
				Value:  toppings.Data,
				Inline: true,
			})
			utils.InteractionRespondMessageEmbed(s, i, embed)
		}
	},
}
