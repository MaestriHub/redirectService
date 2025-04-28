package domain

import (
	"github.com/google/uuid"
)

type Employee struct {
	ID          uuid.UUID
	Name        string
	Description string
}

func NewDummyEmployee(id uuid.UUID) *Employee {
	return &Employee{
		ID:          id,
		Name:        "Employee",
		Description: "It's a dream!",
	}
}
