package roll

import (
	"math/rand"
	"time"
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

		r := rand.NewSource(time.Now().UnixNano())

		min := int(1)
		var max int

		if data.Int("max") == 0 {
			max = int(100)
		} else {
			max = data.Int("max")
		}

		if max < 2 {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Max must be at least 2 for rolling a dice").Build())
		}

		roll := rand.New(r).Intn(max-min+1) + min

		utils.UpdateStatistics("dice_rolled")

		// return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("%d (1-%d)", roll, max).Build())
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitlef("%v rolled a dice", e.User().EffectiveName()).SetDescriptionf("%d (1-%d)", roll, max).Build()).Build())
	}
}
