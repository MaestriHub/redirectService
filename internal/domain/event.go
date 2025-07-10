package domain

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Event interface {
	GetType() string
	GetPayload() ([]byte, error)
}

type EmployeeInviteEvent struct {
	SalonId    uuid.UUID `json:"salonId"`
	EmployeeId uuid.UUID `json:"employeeId"`
}

type SalonInviteEvent struct {
	SalonId uuid.UUID `json:"salonId"`
}

type ClientInviteEvent struct {
	ClientId uuid.UUID `json:"clientId"`
	SalonId  uuid.UUID `json:"salonId"`
}

func (e EmployeeInviteEvent) GetType() string {
	return "EmployeeInvite"
}

func (e EmployeeInviteEvent) GetPayload() ([]byte, error) {
	return json.Marshal(&e)
}

func (e SalonInviteEvent) GetType() string {
	return "SalonInvite"
}

func (e SalonInviteEvent) GetPayload() ([]byte, error) {
	return json.Marshal(&e)
}

func (e ClientInviteEvent) GetType() string {
	return "ClientInvite"
}

func (e ClientInviteEvent) GetPayload() ([]byte, error) {
	return json.Marshal(&e)
}

func NewEvent(eventType string, payload []byte) (Event, error) {
	switch eventType {
	case "EmployeeInvite":
		var e EmployeeInviteEvent
		err := json.Unmarshal(payload, &e)
		if err != nil {
			return nil, err
		}
		return e, nil
	case "SalonInvite":
		var e SalonInviteEvent
		err := json.Unmarshal(payload, &e)
		if err != nil {
			return nil, err
		}
		return e, nil
	case "ClientInvite":
		var e ClientInviteEvent
		err := json.Unmarshal(payload, &e)
		if err != nil {
			return nil, err
		}
		return e, nil
	default:
		return nil, fmt.Errorf("unknown event type %s", eventType)
	}
}

func NewEmployeeInviteEvent(salonId, employeeId uuid.UUID) *EmployeeInviteEvent {
	return &EmployeeInviteEvent{
		SalonId:    salonId,
		EmployeeId: employeeId,
	}
}

func NewSalonInviteEvent(salonId uuid.UUID) *SalonInviteEvent {
	return &SalonInviteEvent{
		SalonId: salonId,
	}
}

func NewClientInviteEvent(clientId, salonId uuid.UUID) *ClientInviteEvent {
	return &ClientInviteEvent{
		ClientId: clientId,
		SalonId:  salonId,
	}
}
