package domain

import (
	"github.com/google/uuid"
)

type Salon struct {
	ID          uuid.UUID
	Name        string
	Description string
}

func NewDummySalon(id uuid.UUID) *Salon {
	return &Salon{
		ID:          id,
		Name:        "Salon",
		Description: "It's a dream!",
	}
}
