package models

import "time"

type AddReminderRequest struct {
	UserId int       `json:"user_id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
}
