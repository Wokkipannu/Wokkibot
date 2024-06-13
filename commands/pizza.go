package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type ToppingsResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

var pizzaCommand = discord.SlashCommandCreate{
	Name:        "pizza",
	Description: "Get random pizza toppings",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionInt{
			Name:        "amount",
			Description: "Amount of toppings",
			Required:    false,
		},
	},
}

func HandlePizza(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		var amount int

		if data.Int("amount") == 0 {
			amount = int(4)
		} else {
			amount = data.Int("amount")
		}

		if amount < 1 {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Choose at least 1 topping").Build())
		}

		res, err := b.Client.Rest().HTTPClient().Get("https://api.wokki.dev/api/random?amount=" + strconv.Itoa(amount))
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while fetching data").Build())
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while reading data").Build())
		}

		toppings := ToppingsResponse{}
		err = json.Unmarshal(body, &toppings)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while parsing data").Build())
		}

		embed := discord.NewEmbedBuilder()
		embed.Author = &discord.EmbedAuthor{
			Name:    fmt.Sprintf("%v's pizza", e.User().EffectiveName()),
			IconURL: *e.User().AvatarURL(),
		}
		runeLength := []rune(toppings.Data)
		var outputs []string
		if len(runeLength) > 900 {
			outputs = append(outputs, string(runeLength[:900]))
			outputs = append(outputs, string(runeLength[900:]))
			embed.AddField("Toppings (1/2)", outputs[0], true)
			embed.AddField("Toppings (2/2)", outputs[1], true)
		} else {
			embed.AddField("Toppings", toppings.Data, false)
		}
		builtEmbed := embed.Build()

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(builtEmbed).Build())
	}
}
