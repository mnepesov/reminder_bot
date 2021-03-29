package models

import (
	"encoding/json"
	"time"
)

type User struct {
	Id       int    `json:"tg_id" binding:"required"`
	ChatId     int64  `json:"chat_id" binding:"required"`
	Username   string `json:"username" binding:"required"`
	FullName   string `json:"full_name" binding:"required"`
	Timezone   string `json:"timezone"`
	RegisterAt time.Time
}

func (u *User) Marshal() ([]byte, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type GetUser struct {
	Id int `json:"id" binding:"required"`
}

type UpdateTimezone struct {
	UserId   int    `json:"user_id"`
	Timezone string `json:"timezone"`
}

type Response struct {
	IsSuccess bool   `json:"is_success"`
	ErrText   string `json:"err_text"`
	Body      []byte `json:"body"`
}

func (resp *Response) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, resp)
	return err
}
