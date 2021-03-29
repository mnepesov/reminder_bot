package models

type NotifyRequest struct {
	Id     int    `json:"id"`
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}
