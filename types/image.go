package types

import (
	"github.com/google/uuid"
)

type Image struct {
	ID     int64     `json:"id"`
	Name   uuid.UUID `json:"name"`
	UserID int       `json:"user_id"`
	Format string    `json:"format"`
}
