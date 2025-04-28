package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBaseInviteEvent_GetType(t *testing.T) {
	event := BaseInviteEvent{
		Type: EmployeeInvite,
	}

	eventType := event.GetType()

	assert.Equal(t, EmployeeInvite, eventType)
}

func TestNewEmployeeInviteEvent(t *testing.T) {
	salonId := uuid.New()
	employeeId := uuid.New()

	event := NewEmployeeInviteEvent(salonId, employeeId)

	assert.NotNil(t, event)
	assert.Equal(t, EmployeeInvite, event.Type)
	assert.Equal(t, salonId, event.SalonId)
	assert.Equal(t, employeeId, event.EmployeeId)
}

func TestNewSalonInviteEvent(t *testing.T) {
	salonId := uuid.New()

	event := NewSalonInviteEvent(salonId)

	assert.NotNil(t, event)
	assert.Equal(t, SalonInvite, event.Type)
	assert.Equal(t, salonId, event.SalonId)
}

func TestNewClientInviteEvent(t *testing.T) {
	clientId := uuid.New()
	employeeId := uuid.New()

	event := NewClientInviteEvent(clientId, employeeId)

	assert.NotNil(t, event)
	assert.Equal(t, ClientInvite, event.Type)
	assert.Equal(t, clientId, event.ClientId)
	assert.Equal(t, employeeId, event.EmployeeId)
}

func TestInviteEvent_InterfaceImplementation(t *testing.T) {
	salonId := uuid.New()
	employeeId := uuid.New()
	clientId := uuid.New()

	employeeEvent := NewEmployeeInviteEvent(salonId, employeeId)
	salonEvent := NewSalonInviteEvent(salonId)
	clientEvent := NewClientInviteEvent(clientId, employeeId)

	var inviteEvent InviteEvent

	inviteEvent = employeeEvent
	assert.Equal(t, EmployeeInvite, inviteEvent.GetType())

	inviteEvent = salonEvent
	assert.Equal(t, SalonInvite, inviteEvent.GetType())

	inviteEvent = clientEvent
	assert.Equal(t, ClientInvite, inviteEvent.GetType())
}

func TestUUID_Initialization(t *testing.T) {
	salonId := uuid.New()
	employeeId := uuid.New()
	clientId := uuid.New()

	employeeEvent := NewEmployeeInviteEvent(salonId, employeeId)
	salonEvent := NewSalonInviteEvent(salonId)
	clientEvent := NewClientInviteEvent(clientId, employeeId)

	assert.NotEqual(t, uuid.Nil, employeeEvent.SalonId)
	assert.NotEqual(t, uuid.Nil, employeeEvent.EmployeeId)

	assert.NotEqual(t, uuid.Nil, salonEvent.SalonId)

	assert.NotEqual(t, uuid.Nil, clientEvent.ClientId)
	assert.NotEqual(t, uuid.Nil, clientEvent.EmployeeId)
}
