package workers

import (
	"log/slog"
	"time"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
)

type Worker struct {
	Bot *wokkibot.Wokkibot
}

func NewWorker(bot *wokkibot.Wokkibot) *Worker {
	return &Worker{
		Bot: bot,
	}
}

func (w *Worker) Start() {
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for range ticker.C {
			w.checkReminders()
		}
	}()
}

func (w *Worker) checkReminders() {
	for _, reminder := range w.Bot.Handlers.ReminderHandler.Reminders {
		if reminder.RemindAt.Before(time.Now()) {
			_, err := w.Bot.Client.Rest().CreateMessage(reminder.ChannelID, discord.NewMessageCreateBuilder().
				SetContentf("<@%s> %s", reminder.UserID.String(), reminder.Message).
				Build())

			if err != nil {
				slog.Error("Error sending reminder", "error", err)
				continue
			}

			w.Bot.Handlers.ReminderHandler.RemoveReminder(reminder.ID)
		}
	}
}
