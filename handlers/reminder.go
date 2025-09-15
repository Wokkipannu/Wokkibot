package handlers

import (
	"sync"
	"time"
	"wokkibot/database"
	"wokkibot/types"

	"github.com/disgoorg/snowflake/v2"
)

type ReminderHandler struct {
	Reminders []types.Reminder
	mu        sync.RWMutex
	updateCh  chan struct{}
}

func NewReminderHandler() *ReminderHandler {
	return &ReminderHandler{
		Reminders: []types.Reminder{},
		updateCh:  make(chan struct{}, 1),
	}
}

func (h *ReminderHandler) LoadReminders() ([]types.Reminder, error) {
	db := database.GetDB()

	rows, err := db.Query("SELECT id, user_id, channel_id, message, remind_at FROM reminders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []types.Reminder
	for rows.Next() {
		var reminder types.Reminder
		err := rows.Scan(&reminder.ID, &reminder.UserID, &reminder.ChannelID, &reminder.Message, &reminder.RemindAt)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}
	return reminders, nil
}

func (h *ReminderHandler) SetReminders(reminders []types.Reminder) {
	h.mu.Lock()
	h.Reminders = reminders
	h.mu.Unlock()
	h.notifyUpdate()
}

func (h *ReminderHandler) UpdateChan() <-chan struct{} { return h.updateCh }

func (h *ReminderHandler) notifyUpdate() {
	select {
	case h.updateCh <- struct{}{}:
	default:
	}
}

func (h *ReminderHandler) AddReminder(reminder types.Reminder) error {
	db := database.GetDB()

	result, err := db.Exec("INSERT INTO reminders (user_id, channel_id, guild_id, message, remind_at) VALUES (?, ?, ?, ?, ?)",
		reminder.UserID, reminder.ChannelID, reminder.GuildID, reminder.Message, reminder.RemindAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	reminder.ID = int(id)
	h.mu.Lock()
	h.Reminders = append(h.Reminders, reminder)
	h.mu.Unlock()

	h.notifyUpdate()

	return nil
}

func (h *ReminderHandler) RemoveReminder(id int) error {
	db := database.GetDB()
	result, err := db.Exec("DELETE FROM reminders WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return nil
	}

	h.mu.Lock()
	newReminders := make([]types.Reminder, 0, len(h.Reminders)-1)
	for _, reminder := range h.Reminders {
		if reminder.ID != id {
			newReminders = append(newReminders, reminder)
		}
	}
	h.Reminders = newReminders
	h.mu.Unlock()

	h.notifyUpdate()

	return nil
}

func (h *ReminderHandler) GetRemindersByUserID(userID snowflake.ID) ([]types.Reminder, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var reminders []types.Reminder
	for _, reminder := range h.Reminders {
		if reminder.UserID == userID {
			reminders = append(reminders, reminder)
		}
	}
	return reminders, nil
}

func (h *ReminderHandler) GetNextRemindAt() (time.Time, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var (
		minTime time.Time
		ok      bool
	)
	for _, r := range h.Reminders {
		if !ok || r.RemindAt.Before(minTime) {
			minTime = r.RemindAt
			ok = true
		}
	}
	return minTime, ok
}

func (h *ReminderHandler) GetDueReminders(now time.Time) []types.Reminder {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var due []types.Reminder
	for _, r := range h.Reminders {
		if !r.RemindAt.After(now) {
			due = append(due, r)
		}
	}
	return due
}
