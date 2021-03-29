package notify

import "reminder_bot/internal/models"

type Repository interface {
	GetNotifies() ([]models.NotifyRequest, error)
	DeactivateReminder(id int) error
}