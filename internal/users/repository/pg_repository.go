package repository

import (
	"github.com/jmoiron/sqlx"
	"reminder_bot/internal/models"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

func (u *UserPostgres) GetUserById(tgId int) (models.User, error) {
	user := models.User{}
	query := "SELECT id, chat_id, username, full_name, timezone FROM users WHERE id = $1"
	err := u.db.QueryRow(query, tgId).Scan(&user.Id, &user.ChatId, &user.Username, &user.FullName, &user.Timezone)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *UserPostgres) CreateUser(user models.User) error {
	query := "INSERT INTO users(id, chat_id, username, full_name, timezone) values ($1,$2,$3,$4,'europe/minsk')"
	_, err := u.db.Exec(query, user.Id, user.ChatId, user.Username, user.FullName)
	return err
}

func (u *UserPostgres) UpdateTimezone(userId int, timezone string) error {
	query := "UPDATE users SET timezone = $1 WHERE id = $2"
	_, err := u.db.Exec(query, timezone, userId)
	return err
}
