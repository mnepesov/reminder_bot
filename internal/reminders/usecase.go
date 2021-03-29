package reminders

import "reminder_bot/internal/models"

type UseCase interface {
	AddReminders(reminder models.AddReminderRequest) error
}
