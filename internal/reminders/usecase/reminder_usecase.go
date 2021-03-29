package usecase

import (
	"reminder_bot/internal/models"
	"reminder_bot/internal/reminders"
)

type ReminderUseCase struct {
	repo reminders.Repository
}

func NewReminderUseCase(repo reminders.Repository) *ReminderUseCase {
	return &ReminderUseCase{
		repo: repo,
	}
}

func (r *ReminderUseCase) AddReminders(reminder models.AddReminderRequest) error {
	return r.repo.AddReminders(reminder)
}