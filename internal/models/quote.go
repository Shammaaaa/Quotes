package models

import "time"

type Quote struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`
	Text      string    `json:"quote"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRequest struct {
	Author string `json:"author" validate:"required"`
	Text   string `json:"quote" validate:"required"`
}
