package notify

import "reminder_bot/internal/models"

type UseCase interface {
	GetNotifies() ([]models.NotifyRequest, error)
	DeactivateReminder(id int) error
}
