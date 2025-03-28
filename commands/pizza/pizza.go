package pizza

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"wokkibot/database"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type Topping struct {
	Name  string
	Count int64
}

var PizzaCommand = discord.SlashCommandCreate{
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

const (
	MAX_FIELD_LENGTH = 1024
)

func formatNumber(n int64) string {
	str := fmt.Sprintf("%d", n)

	var result strings.Builder
	length := len(str)
	for i, char := range str {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteRune(' ')
		}
		result.WriteRune(char)
	}
	return result.String()
}

func HandlePizza(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		data := e.SlashCommandInteractionData()

		var amount int

		if data.Int("amount") == 0 {
			amount = int(4)
		} else {
			amount = data.Int("amount")
		}

		if amount < 1 {
			utils.HandleError(e, "Choose at least 1 topping", "pizza")
			return nil
		}

		output, err := getRandomToppings(amount)
		if err != nil {
			utils.HandleError(e, err.Error(), "pizza")
			return err
		}

		embed := discord.NewEmbedBuilder()
		embed.Author = &discord.EmbedAuthor{
			Name:    fmt.Sprintf("%v's pizza", e.User().EffectiveName()),
			IconURL: *e.User().AvatarURL(),
		}
		embed.SetColor(utils.COLOR_BLURPLE)

		var currentChunk strings.Builder
		var chunks []string

		for i, tc := range output {
			toppingText := tc

			if currentChunk.Len()+len(toppingText)+2 > MAX_FIELD_LENGTH {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
			}

			if currentChunk.Len() > 0 {
				currentChunk.WriteString(", ")
			}
			currentChunk.WriteString(toppingText)

			if i == len(output)-1 && currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
			}
		}

		if len(chunks) > 1 {
			for i, chunk := range chunks {
				embed.AddField(fmt.Sprintf("Toppings (x%d) (%d/%d)", amount, i+1, len(chunks)), chunk, false)
			}
		} else if len(chunks) == 1 {
			embed.AddField(fmt.Sprintf("Toppings (x%d)", amount), chunks[0], false)
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(embed.Build()).
			AddActionRow(discord.NewPrimaryButton("Rerandomize", "/pizza/randomize").WithEmoji(discord.ComponentEmoji{Name: "🍕"})).
			Build())

		utils.UpdateStatistics("pizzas_generated")

		return err
	}
}

func HandlePizzaRandomize(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		message := e.Message
		if len(message.Embeds) == 0 {
			return errors.New("no embeds found in message")
		}

		embedAuthor := message.Embeds[0].Author
		if embedAuthor == nil {
			return errors.New("no author found in embed")
		}

		author := message.InteractionMetadata.User
		if author.ID != e.User().ID {
			return e.Respond(discord.InteractionResponseTypeCreateMessage, discord.NewMessageCreateBuilder().
				SetContent("Only the original user can rerandomize the pizza!").
				SetEphemeral(true).
				Build())
		}

		titleField := message.Embeds[0].Fields[0]
		title := titleField.Name

		var amount int
		_, err := fmt.Sscanf(title, "Toppings (x%d)", &amount)
		if err != nil {
			amount = 4
		}

		output, err := getRandomToppings(amount)
		if err != nil {
			return err
		}

		embed := discord.NewEmbedBuilder()
		embed.Author = &discord.EmbedAuthor{
			Name:    fmt.Sprintf("%v's pizza", e.User().EffectiveName()),
			IconURL: *e.User().AvatarURL(),
		}
		embed.SetColor(utils.COLOR_BLURPLE)

		var currentChunk strings.Builder
		var chunks []string

		for i, tc := range output {
			toppingText := tc

			if currentChunk.Len()+len(toppingText)+2 > MAX_FIELD_LENGTH {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
			}

			if currentChunk.Len() > 0 {
				currentChunk.WriteString(", ")
			}
			currentChunk.WriteString(toppingText)

			if i == len(output)-1 && currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
			}
		}

		if len(chunks) > 1 {
			for i, chunk := range chunks {
				embed.AddField(fmt.Sprintf("Toppings (x%d) (%d/%d)", amount, i+1, len(chunks)), chunk, false)
			}
		} else if len(chunks) == 1 {
			embed.AddField(fmt.Sprintf("Toppings (x%d)", amount), chunks[0], false)
		}

		err = e.Respond(discord.InteractionResponseTypeUpdateMessage, discord.NewMessageUpdateBuilder().
			SetEmbeds(embed.Build()).
			AddActionRow(discord.NewPrimaryButton("Rerandomize", "/pizza/randomize").WithEmoji(discord.ComponentEmoji{Name: "🍕"})).
			Build())

		if err != nil {
			return err
		}

		return nil
	}
}

func getRandomToppings(amount int) ([]string, error) {
	db := database.GetDB()

	rows, err := db.Query("SELECT name FROM pizza_toppings")
	if err != nil {
		return nil, err
	}

	var allToppings []string
	for rows.Next() {
		var topping string
		if err := rows.Scan(&topping); err != nil {
			return nil, err
		}
		allToppings = append(allToppings, topping)
	}

	if len(allToppings) == 0 {
		return nil, errors.New("no toppings found in database")
	}

	toppingsCount := len(allToppings)

	randomToppings := make([]Topping, toppingsCount)
	amountToRoll := int64(amount)

	for amountToRoll > 0 {
		roll := rand.Intn(toppingsCount)
		countLimit := math.Ceil(float64(amount) / 12)
		count := rand.Int63n(int64(countLimit))
		if count == 0 {
			count += 1
		}

		randomToppings[roll].Count += count
		randomToppings[roll].Name = allToppings[roll]
		amountToRoll -= count
	}

	var output []string
	for _, v := range randomToppings {
		if v.Count > 0 {
			if v.Count > 1 {
				output = append(output, fmt.Sprintf("%sx %s", formatNumber(v.Count), v.Name))
			} else {
				output = append(output, v.Name)
			}
		}
	}

	return output, nil
}
