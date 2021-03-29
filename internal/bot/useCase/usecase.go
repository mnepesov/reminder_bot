package useCase

import "reminder_bot/internal/models"

type UserUseCase interface {
	GetUser(userId int) (models.User, error)
	CreateUser(user models.User) error
	UpdateTimezone(userId int, timezone string) error
}

type ParsingService interface {
	Parse(req models.ParseRequest) (models.ParseResponse, error)
}

type ReminderService interface {
	AddReminder(req models.AddReminderRequest) error
}

type UseCase interface {
	UserUseCase
	ParsingService
	ReminderService
}
