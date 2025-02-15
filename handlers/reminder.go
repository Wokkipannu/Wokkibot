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

	result, err := db.Exec("INSERT INTO reminders (user_id, channel_id, message, remind_at) VALUES (?, ?, ?, ?)",
		reminder.UserID, reminder.ChannelID, reminder.Message, reminder.RemindAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	reminder.ID = int(id)
	h.Reminders = append(h.Reminders, reminder)

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

	newReminders := make([]types.Reminder, 0, len(h.Reminders)-1)
	for _, reminder := range h.Reminders {
		if reminder.ID != id {
			newReminders = append(newReminders, reminder)
		}
	}
	h.Reminders = newReminders

	return nil
}
