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

var SetReminderCommand = discord.ApplicationCommandOptionSubCommand{
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
