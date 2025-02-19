package payload

import uuid "github.com/google/uuid"

type Employeer struct {
	ID      uuid.UUID
	SalonId uuid.UUID
}

type Salon struct {
	ID uuid.UUID
}

type Customer struct {
	ID uuid.UUID
}

type MasterToSalon struct {
	EmployeeId uuid.UUID
}
