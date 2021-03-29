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

func (r *ReminderPostgres) GetRemindersByUserId(userId int) ([]models.Reminder, error) {
	var reminders []models.Reminder
	rows, err := r.db.Query("select id, user_id, text, date from reminders where user_id = $1 and is_active = true ", userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		reminder := models.Reminder{}
		err := rows.Scan(&reminder.Id, &reminder.UserId, &reminder.Text, &reminder.Date)
		if err != nil {
			continue
		}

		reminders = append(reminders, reminder)
	}

	return reminders, nil

}
