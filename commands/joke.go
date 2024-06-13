package commands

import (
	"encoding/json"
	"io"
	"log/slog"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type JokeRes struct {
	Error    bool   `json:"error"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Joke     string `json:"joke,omitempty"`
	Setup    string `json:"setup,omitempty"`
	Delivery string `json:"delivery,omitempty"`
	Flags    struct {
		NSFW      bool `json:"nsfw"`
		Religious bool `json:"religious"`
		Political bool `json:"political"`
		Racist    bool `json:"racist"`
		Sexist    bool `json:"sexist"`
		Explicit  bool `json:"explicit"`
	} `json:"flags"`
	ID   int    `json:"id"`
	Safe bool   `json:"safe"`
	Lang string `json:"lang"`
}

var jokeCommand = discord.SlashCommandCreate{
	Name:        "joke",
	Description: "Replies with a random joke",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "category",
			Description: "Category of the joke",
			Required:    false,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "Any",
					Value: "Any",
				},
				{
					Name:  "Programming",
					Value: "Programming",
				},
				{
					Name:  "Miscellaneous",
					Value: "Miscellaneous",
				},
				{
					Name:  "Dark",
					Value: "Dark",
				},
				{
					Name:  "Pun",
					Value: "Pun",
				},
				{
					Name:  "Spooky",
					Value: "Spooky",
				},
				{
					Name:  "Christmas",
					Value: "Christmas",
				},
			},
		},
	},
}

func HandleJoke(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		var category string
		if c, ok := data.OptString("category"); ok {
			category = c
		} else {
			category = "Any"
		}

		res, err := b.Client.Rest().HTTPClient().Get("https://v2.jokeapi.dev/joke/" + category)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while fetching joke").Build())
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while reading joke").Build())
		}

		var jokeRes JokeRes
		err = json.Unmarshal(body, &jokeRes)
		if err != nil {
			slog.Error("Error while unmarshaling joke response", slog.Any("err", err))
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while parsing joke").Build())
		}

		embed := discord.NewEmbedBuilder()

		if jokeRes.Joke == "" {
			embed.SetDescriptionf("%v\n\n%v", jokeRes.Setup, jokeRes.Delivery)
		} else {
			embed.SetDescription(jokeRes.Joke)
		}

		embed.SetColor(utils.RGBToInteger(255, 215, 0))
		embed.SetFooterTextf("ID %v", jokeRes.ID)

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
	}
}
