package params

import (
	"github.com/google/uuid"
)

type CreateEmployeeInviteLink struct {
	EmployeeId uuid.UUID `json:"employeeId"`
	SalonId    uuid.UUID `json:"salonId"`
}

type CreateSalonInviteLink struct {
	SalonId uuid.UUID `json:"salonId"`
}

type CreateClientInviteLink struct {
	EmployeeId uuid.UUID `json:"employeeId"`
	SalonId    uuid.UUID `json:"salonId"`
}
