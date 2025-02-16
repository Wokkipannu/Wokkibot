package remind

import (
	"fmt"
	"time"
	"wokkibot/types"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var RemindCommand = discord.SlashCommandCreate{
	Name:        "remind",
	Description: "Remind you to do something at a specific time",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "set",
			Description: "Set a reminder",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "reminder",
					Description: "The reminder message",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "hours",
					Description: "The amount of hours after you want to be reminded",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "minutes",
					Description: "The amount of minutes after you want to be reminded",
					Required:    true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "list",
			Description: "List your reminders",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "delete",
			Description: "Delete a reminder",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "id",
					Description: "The ID of the reminder to delete",
					Required:    true,
				},
			},
		},
	},
}

func HandleRemind(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		data := e.SlashCommandInteractionData()
		reminder := data.String("reminder")
		hours := data.Int("hours")
		minutes := data.Int("minutes")

		remindAt := time.Now().Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute)

		b.Handlers.ReminderHandler.AddReminder(types.Reminder{
			UserID:    e.User().ID,
			ChannelID: e.Channel().ID(),
			GuildID:   *e.GuildID(),
			Message:   reminder,
			RemindAt:  remindAt,
		})

		embed := discord.NewEmbedBuilder().
			SetTitle("Reminder set").
			SetColor(utils.COLOR_BLURPLE).
			AddField("Reminder", reminder, true).
			AddField("Time", fmt.Sprintf("<t:%d:R>", remindAt.Unix()), true).
			Build()

		_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(embed).
			Build())

		return err
	}
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
