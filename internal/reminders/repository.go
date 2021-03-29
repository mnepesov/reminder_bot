package reminders

import "reminder_bot/internal/models"

type Repository interface {
	AddReminders(reminder models.AddReminderRequest) error
	GetRemindersByUserId(userId int) ([]models.Reminder, error)
}

