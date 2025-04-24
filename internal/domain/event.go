package domain

import "github.com/google/uuid"

type Event = string

const (
	EmployeeInvite Event = "EmployeeInvite"
	SalonInvite    Event = "SalonInvite"
	ClientInvite   Event = "ClientInvite"
)

type InviteEvent interface {
	GetType() string
}

type BaseInviteEvent struct {
	Type Event `json:"type"`
}

func (e BaseInviteEvent) GetType() string {
	return e.Type
}

type EmployeeInviteEvent struct {
	BaseInviteEvent
	SalonId    uuid.UUID `json:"salonId"`
	EmployeeId uuid.UUID `json:"employeeId"`
}

type SalonInviteEvent struct {
	BaseInviteEvent
	SalonId uuid.UUID `json:"salonId"`
}

type ClientInviteEvent struct {
	BaseInviteEvent
	ClientId   uuid.UUID `json:"clientId"`
	EmployeeId uuid.UUID `json:"employeeId"`
}

func NewEmployeeInviteEvent(salonId, employeeId uuid.UUID) *EmployeeInviteEvent {
	return &EmployeeInviteEvent{
		BaseInviteEvent: BaseInviteEvent{
			Type: EmployeeInvite,
		},
		SalonId:    salonId,
		EmployeeId: employeeId,
	}
}

func NewSalonInviteEvent(salonId uuid.UUID) *SalonInviteEvent {
	return &SalonInviteEvent{
		BaseInviteEvent: BaseInviteEvent{
			Type: SalonInvite,
		},
		SalonId: salonId,
	}
}

func NewClientInviteEvent(clientId, employeeId uuid.UUID) *ClientInviteEvent {
	return &ClientInviteEvent{
		BaseInviteEvent: BaseInviteEvent{
			Type: ClientInvite,
		},
		ClientId:   clientId,
		EmployeeId: employeeId,
	}
}
