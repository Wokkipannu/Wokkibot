package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
	"wokkibot/config"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

type Topping struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Topping   string `json:"topping"`
}

type ToppingsResponse struct {
	Data    []Topping `json:"data"`
	Message string    `json:"message"`
	Status  string    `json:"status"`
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
		res, err := s.Client.Get(API)
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

		selectedToppings := make(map[string]int)

		for i := 0; i < amount; i++ {
			selectedToppings[toppings.Data[rand.Intn(len(toppings.Data))].Topping] += 1
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
			Name:    fmt.Sprintf("%vn pizza", i.Member.Nick),
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

		var output []string
		for k, v := range selectedToppings {
			if v > 1 {
				output = append(output, fmt.Sprintf("%vx %v", v, k))
			} else {
				output = append(output, k)
			}
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Täytteet",
			Value:  fmt.Sprintf("%v", strings.Join(output[:], ", ")),
			Inline: true,
		})

		utils.InteractionRespondMessageEmbed(s, i, embed)
	},
}
