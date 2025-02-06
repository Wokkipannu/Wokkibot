package friday

import (
	"wokkibot/database"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var FridayCommand = discord.SlashCommandCreate{
	Name:        "friday",
	Description: "Post a friday celebration video",
}

func HandleFriday(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		db := database.GetDB()

		var url string
		err := db.QueryRow("SELECT url FROM friday_clips ORDER BY RANDOM() LIMIT 1").Scan(&url)
		if err != nil {
			e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("No friday clips found").
				Build(),
			)
			return err
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(url).
			Build(),
		)
	}
}
