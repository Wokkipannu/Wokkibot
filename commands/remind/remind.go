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
			Message:   reminder,
			RemindAt:  remindAt,
		})

		embed := discord.NewEmbedBuilder().
			SetTitle("Reminder set").
			SetColor(utils.COLOR_BLURPLE).
			SetDescription("I'll remind you about: \"" + reminder + "\" in " + formatDuration(int64(hours), int64(minutes))).
			Build()

		_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(embed).
			Build())

		return err
	}
}

func formatDuration(hours, minutes int64) string {
	var result string
	if hours > 0 {
		if hours == 1 {
			result += "1 hour"
		} else {
			result += fmt.Sprintf("%d hours", hours)
		}
	}
	if minutes > 0 {
		if hours > 0 {
			result += " and "
		}
		if minutes == 1 {
			result += "1 minute"
		} else {
			result += fmt.Sprintf("%d minutes", minutes)
		}
	}
	return result
}
