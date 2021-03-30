package usecase

import (
	"reminder_bot/internal/models"
	"reminder_bot/internal/reminder"
)

type ReminderUseCase struct {
	repo reminder.Repository
}

func NewReminderUseCase(repo reminder.Repository) *ReminderUseCase {
	return &ReminderUseCase{
		repo: repo,
	}
}

func (r *ReminderUseCase) AddReminders(reminder models.AddReminderRequest) error {
	return r.repo.AddReminders(reminder)
}
func (r *ReminderUseCase) GetRemindersByUserId(userId int) ([]models.Reminder, error) {
	return r.repo.GetRemindersByUserId(userId)
}
