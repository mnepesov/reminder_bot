package users

import "reminder_bot/internal/models"

type UseCase interface {
	GetUserById(tgId int) ([]byte, error)
	CreateUser(user models.User) error
	UpdateTimezone(userId int, timezone string) error
}
