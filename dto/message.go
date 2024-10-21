package dto

import "time"

type MessagesResponse struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Sender     string    `json:"sender"`      // Sender's username
	Receiver   string    `json:"receiver"`    // Receiver's username
	Content    string    `json:"content"`     // Message content
	IsAccepted bool      `json:"is_accepted"` // Message acceptance status
}
