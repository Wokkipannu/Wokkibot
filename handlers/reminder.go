package handlers

import (
	"wokkibot/database"
	"wokkibot/types"
)

type ReminderHandler struct {
	Reminders []types.Reminder
}

func NewReminderHandler() *ReminderHandler {
	return &ReminderHandler{
		Reminders: []types.Reminder{},
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

func (h *ReminderHandler) AddReminder(reminder types.Reminder) error {
	db := database.GetDB()

	_, err := db.Exec("INSERT INTO reminders (user_id, channel_id, message, remind_at) VALUES (?, ?, ?, ?)", reminder.UserID, reminder.ChannelID, reminder.Message, reminder.RemindAt)

	h.Reminders = append(h.Reminders, reminder)

	return err
}

func (h *ReminderHandler) RemoveReminder(id int) error {
	db := database.GetDB()

	_, err := db.Exec("DELETE FROM reminders WHERE id = ?", id)

	for i, reminder := range h.Reminders {
		if reminder.ID == id {
			h.Reminders = append(h.Reminders[:i], h.Reminders[i+1:]...)
			break
		}
	}

	return err
}
