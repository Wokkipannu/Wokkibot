package remind

import (
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var DeleteReminderCommand = discord.ApplicationCommandOptionSubCommand{
	Name:        "delete",
	Description: "Delete a reminder",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionInt{
			Name:        "id",
			Description: "The ID of the reminder to delete",
			Required:    true,
		},
	},
}

func HandleDeleteMyReminder(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		reminderID := int(e.SlashCommandInteractionData().Int("id"))

		reminders, err := b.Handlers.ReminderHandler.GetRemindersByUserID(e.User().ID)
		if err != nil {
			return err
		}

		var found bool
		for _, reminder := range reminders {
			if reminder.ID == reminderID {
				found = true
				break
			}
		}

		if !found {
			_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetContent("You don't have a reminder with that ID.").
				Build())
			return err
		}

		err = b.Handlers.ReminderHandler.RemoveReminder(reminderID)
		if err != nil {
			return err
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent("Reminder deleted successfully.").
			Build())

		return err
	}
}
