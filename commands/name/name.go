package name

import (
	"fmt"
	"wokkibot/database"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var NameCommand = discord.SlashCommandCreate{
	Name:        "name",
	Description: "Generates a random two-part name from the names list",
}

func Init() {}

func HandleName(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		db := database.GetDB()

		rows, err := db.Query("SELECT name FROM names ORDER BY RANDOM() LIMIT 2")
		if err != nil {
			utils.HandleError(e, "Error while fetching names", err.Error())
			return err
		}
		defer rows.Close()

		var names []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				utils.HandleError(e, "Error while scanning names", err.Error())
				return err
			}
			names = append(names, name)
		}

		if len(names) < 2 {
			utils.HandleError(e, "Not enough names in the database", "")
			return fmt.Errorf("not enough names in database")
		}

		randomName := fmt.Sprintf("You are **%s%s**", names[0], names[1])

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent(randomName).
			Build())

		utils.UpdateStatistics("names_given")

		return err
	}
}
