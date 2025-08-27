package models

import "time"

type Event struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Element   string    `json:"element"`
	Timestamp time.Time `json:"timestamp"`
	Duration  float64   `json:"duration"`
}