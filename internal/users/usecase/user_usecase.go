package usecase

import (
	"reminder_bot/internal/models"
	"reminder_bot/internal/users"
)

type UserUseCase struct {
	repo users.Repository
}

func NewUserUseCase(repo users.Repository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (u *UserUseCase) GetUserById(tgId int) ([]byte, error) {
	user, err := u.repo.GetUserById(tgId)
	if err != nil {
		return nil, err
	}

	return user.Marshal()
}

func (u *UserUseCase) CreateUser(user models.User) error {
	return u.repo.CreateUser(user)
}

func (u *UserUseCase) UpdateTimezone(userId int, timezone string) error {
	return u.repo.UpdateTimezone(userId, timezone)
}
