package flip

import (
	"math/rand"
	"time"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var FlipCommand = discord.SlashCommandCreate{
	Name:        "flip",
	Description: "Flip a coin",
}

func HandleFlip(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		r := rand.NewSource(time.Now().UnixNano())
		min := int(0)
		max := int(1)

		flip := rand.New(r).Intn(max-min+1) + min

		var result string
		if flip == 0 {
			result = "Heads"
		} else {
			result = "Tails"
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(result).Build())
	}
}
