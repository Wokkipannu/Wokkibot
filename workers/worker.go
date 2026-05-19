package workers

import (
	"log/slog"
	"time"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
)

const ytdlpUpdateInterval = 10 * time.Minute

type Worker struct {
	Bot *wokkibot.Wokkibot
}

func NewWorker(bot *wokkibot.Wokkibot) *Worker {
	return &Worker{
		Bot: bot,
	}
}

func (w *Worker) Start() {
	go w.runReminderScheduler()
	go w.runYtdlpUpdater()
}

func (w *Worker) runYtdlpUpdater() {
	w.checkYtdlpUpdate()

	ticker := time.NewTicker(ytdlpUpdateInterval)
	defer ticker.Stop()

	for range ticker.C {
		w.checkYtdlpUpdate()
	}
}

func (w *Worker) checkYtdlpUpdate() {
	current := utils.GetYtdlpVersion()
	latest, err := utils.GetLatestYtdlpVersion()
	if err != nil {
		slog.Error("Error checking latest yt-dlp version", "error", err)
		return
	}

	if current == latest {
		return
	}

	slog.Info("Updating yt-dlp", "current_version", current, "latest_version", latest)
	if err := utils.UpdateYtdlpBinary(); err != nil {
		slog.Error("Error updating yt-dlp", "error", err)
		return
	}

	slog.Info("yt-dlp update completed", "version", latest)
}

func (w *Worker) runReminderScheduler() {
	var timer *time.Timer
	resetTimer := func(d time.Duration) {
		if timer == nil {
			timer = time.NewTimer(d)
			return
		}
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(d)
	}

	recalc := func() {
		if nextAt, ok := w.Bot.Handlers.ReminderHandler.GetNextRemindAt(); ok {
			now := time.Now()
			if nextAt.After(now) {
				resetTimer(nextAt.Sub(now))
			} else {
				resetTimer(0)
			}
		} else {
			resetTimer(24 * time.Hour)
		}
	}

	recalc()
	for {
		select {
		case <-w.Bot.Handlers.ReminderHandler.UpdateChan():
			recalc()
		case <-timer.C:
			w.fireDueReminders()
			recalc()
		}
	}
}

func (w *Worker) fireDueReminders() {
	now := time.Now()
	due := w.Bot.Handlers.ReminderHandler.GetDueReminders(now)
	for _, reminder := range due {
		_, err := w.Bot.Client.Rest.CreateMessage(reminder.ChannelID, discord.NewMessageCreate().
			WithContentf("<@%s> %s", reminder.UserID.String(), reminder.Message))
		if err != nil {
			slog.Error("Error sending reminder", "error", err)
			continue
		}

		if err := w.Bot.Handlers.ReminderHandler.RemoveReminder(reminder.ID); err != nil {
			slog.Error("Error removing reminder", "error", err)
		}
	}
}
