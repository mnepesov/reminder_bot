package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"reminder_bot/internal/models"
)

type NotifyPostgres struct {
	db *sqlx.DB
}

func NewNotifyPostgres(db *sqlx.DB) *NotifyPostgres {
	return &NotifyPostgres{db: db}
}

func (n *NotifyPostgres) GetNotifies() ([]models.NotifyRequest, error) {
	var notifies []models.NotifyRequest
	rows, err := n.db.Query(`
		select r.id,
			   u.chat_id,
			   r.text
		from reminder r
				 join users u on u.id = r.user_id
		where timezone(u.timezone, r.date) between timezone(u.timezone, now()) and timezone(u.timezone, now() + (1 * interval '1 minute'))
	`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		notify := models.NotifyRequest{}
		err := rows.Scan(&notify.Id, &notify.ChatId, &notify.Text)
		if err != nil {
			fmt.Println(err)
			continue
		}
		
		notifies = append(notifies, notify)
		
	}
	return notifies, nil
}

func (r *NotifyPostgres) DeactivateReminder(id int) error {
	_, err := r.db.Exec("update reminder set is_active = false where id = $1", id)
	return err
}
