package roll

import (
	"crypto/rand"
	"math/big"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var RollCommand = discord.SlashCommandCreate{
	Name:        "roll",
	Description: "Roll a dice",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionInt{
			Name:        "max",
			Description: "Highest possible roll value",
			Required:    false,
		},
	},
}

func HandleRoll(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		min := 1
		max := 100

		if data.Int("max") != 0 {
			max = data.Int("max")
		}

		if max < 2 {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Max must be at least 2 for rolling a dice").Build())
		}

		n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to generate random number").Build())
		}
		roll := int(n.Int64()) + min

		utils.UpdateStatistics("dice_rolled")

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitlef("%v rolled a dice", e.User().EffectiveName()).SetDescriptionf("%d (1-%d)", roll, max).Build()).Build())
	}
}
