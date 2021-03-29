package users

import "reminder_bot/internal/models"

type Repository interface {
	GetUserById(tgId int) (models.User, error)
	CreateUser(user models.User) error
	UpdateTimezone(userId int, timezone string) error
}