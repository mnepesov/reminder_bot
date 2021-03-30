package reminder

import "reminder_bot/internal/models"

type UseCase interface {
	AddReminders(reminder models.AddReminderRequest) error
	GetRemindersByUserId(userId int) ([]models.Reminder, error)
}
