package repository

import (
	"github.com/jmoiron/sqlx"
	"reminder_bot/internal/models"
)

type ReminderPostgres struct {
	db *sqlx.DB
}

func NewReminderPostgres(db *sqlx.DB) *ReminderPostgres {
	return &ReminderPostgres{
		db: db,
	}
}

func (r *ReminderPostgres) AddReminders(reminder models.AddReminderRequest) error {
	
	_, err := r.db.Exec("insert into reminders(user_id, text, date) values ($1, $2, $3)", reminder.UserId, reminder.Text, reminder.Date)
	
	return err
}
