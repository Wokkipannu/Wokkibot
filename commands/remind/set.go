package remind

import (
	"fmt"
	"strconv"
	"strings"
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
		discord.ApplicationCommandOptionString{
			Name:        "time",
			Description: "24h time in HH:MM",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "date",
			Description: "Optional date in DD.MM.YYYY",
			Required:    false,
		},
		discord.ApplicationCommandOptionString{
			Name:        "timezone",
			Description: "Optional timezone like Europe/Helsinki",
			Required:    false,
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
		timeStr := data.String("time")
		dateStr, _ := data.OptString("date")
		tzStr, _ := data.OptString("timezone")

		loc := time.Local
		if tzStr != "" {
			if l, err := time.LoadLocation(tzStr); err == nil {
				loc = l
			}
		}

		parts := strings.Split(timeStr, ":")
		if len(parts) != 2 {
			_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetContent("Invalid time format. Use HH:MM, e.g., 15:30").
				Build())
			return err
		}

		hour, herr := parseTwoDigit(parts[0])
		minute, merr := parseTwoDigit(parts[1])
		if herr != nil || merr != nil || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
			_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetContent("Invalid time. Hours 00-23 and minutes 00-59.").
				Build())
			return err
		}

		now := time.Now().In(loc)

		var target time.Time
		if dateStr != "" {
			dp := strings.Split(dateStr, ".")
			if len(dp) != 3 {
				_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
					SetContent("Invalid date format. Use DD.MM.YYYY, e.g., 20.09.2025").
					Build())
				return err
			}
			day, derr1 := strconv.Atoi(dp[0])
			month, derr2 := strconv.Atoi(dp[1])
			year, derr3 := strconv.Atoi(dp[2])
			if derr1 != nil || derr2 != nil || derr3 != nil || day < 1 || day > 31 || month < 1 || month > 12 || year < 1 {
				_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
					SetContent("Invalid date value. Use DD.MM.YYYY").
					Build())
				return err
			}
			target = time.Date(year, time.Month(month), day, hour, minute, 0, 0, loc)
			if !target.After(now) {
				_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
					SetContent("The specified date/time is in the past.").
					Build())
				return err
			}
		} else {
			target = time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, loc)
			if !target.After(now) {
				target = target.Add(24 * time.Hour)
			}
		}

		remindAt := target.In(time.Local)

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

func parseTwoDigit(s string) (int, error) {
	if len(s) == 1 {
		s = "0" + s
	}
	if len(s) != 2 {
		return 0, fmt.Errorf("invalid")
	}
	return strconv.Atoi(s)
}
