package remind

import (
	"fmt"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var ListRemindersCommand = discord.ApplicationCommandOptionSubCommand{
	Name:        "list",
	Description: "List your reminders",
}

func HandleListMyReminders(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		reminders, err := b.Handlers.ReminderHandler.GetRemindersByUserID(e.User().ID)
		if err != nil {
			return err
		}

		if len(reminders) == 0 {
			_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetContent("You have no reminders.").
				Build())

			return err
		}

		embedBuilder := discord.NewEmbedBuilder().
			SetTitle("Your reminders").
			SetColor(utils.COLOR_BLURPLE)

		for _, reminder := range reminders {
			embedBuilder.AddField(
				fmt.Sprintf("ID: %d | %s", reminder.ID, reminder.Message),
				fmt.Sprintf("<t:%d:R>", reminder.RemindAt.Unix()),
				false,
			)
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(embedBuilder.Build()).
			Build())

		return err
	}
}
