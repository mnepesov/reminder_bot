package models

import "time"

type AddReminderRequest struct {
	UserId int       `json:"user_id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
}

type Reminder struct {
	Id       int       `json:"id"`
	UserId   int       `json:"user_id"`
	Text     string    `json:"text"`
	Date     time.Time `json:"date"`
	IsActive bool      `json:"is_active"`
}

type GetRemindersRequest struct {
	UserId int `json:"user_id"`
}
