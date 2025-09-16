package models

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID `json:"id"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
