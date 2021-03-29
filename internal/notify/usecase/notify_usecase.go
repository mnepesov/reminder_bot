package usecase

import (
	"reminder_bot/internal/models"
	"reminder_bot/internal/notify"
)

type NotifyUseCase struct {
	repo notify.Repository
}

func NewNotifyUseCase(repo notify.Repository) *NotifyUseCase {
	return &NotifyUseCase{
		repo: repo,
	}
}

func (n *NotifyUseCase) GetNotifies() ([]models.NotifyRequest, error) {
	return n.repo.GetNotifies()
}
func (n *NotifyUseCase) DeactivateReminder(id int) error {
	return n.repo.DeactivateReminder(id)
}