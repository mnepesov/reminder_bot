package cache

import "reminder_bot/internal/models"

type Cache interface {
	GetUserInfo(userId int) (*models.User, error)
	SetUserInfo(userId int, user *models.User) error
	SetNavigation(userId int, navigation string) error
	GetNavigation(userId int) (string, error)
}
