package models

import "github.com/google/uuid"

type World struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
