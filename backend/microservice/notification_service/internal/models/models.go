package models

import "time"

// Response represents the structure of an API response
type Response struct {
	Errors []string    `json:"Errors"`
	Data   interface{} `json:"Data"`
}

// Notification represents the structure of a notification
type Notification struct {
	ID       int
	UserID   int
	SenderID int
	Type     string
	Message  string
	Time     time.Time
	IsRead   bool
}
