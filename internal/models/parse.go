package models

import "time"

type ParseRequest struct {
	Text string `json:"text"`
}

type ParseResponse struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}
