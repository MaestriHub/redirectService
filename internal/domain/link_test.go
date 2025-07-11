package domain

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDirectLink(t *testing.T) {
	// Arrange
	salonId := uuid.New()
	event := NewSalonInviteEvent(salonId)

	// Act
	link := NewDirectLink(event)

	// Assert
	assert.NotNil(t, link)
	assert.NotEmpty(t, link.NanoId, "NanoId should not be empty")
	assert.Equal(t, 0, link.Clicks, "Clicks should be initialized to 0")
	assert.Equal(t, event.GetType(), link.Event.GetType(), "Event type should match")

	// Check Payload deserialization
	payload, err := link.Event.GetPayload()
	if err != nil {
		t.Fatalf("Failed to get payload: %v", err)
	}

	var payloadEvent SalonInviteEvent
	err = json.Unmarshal(payload, &payloadEvent)
	if err != nil {
		t.Fatalf("Failed to unmarshal payload: %v", err)
	}

	assert.NoError(t, err, "Payload should be valid JSON")
	assert.Equal(t, salonId, payloadEvent.SalonId, "Payload should contain correct SalonId")
}

func TestDirectLink_IncClicks(t *testing.T) {
	// Arrange
	link := DirectLink{Clicks: 5}

	// Act
	link.IncClicks()

	// Assert
	assert.Equal(t, 6, link.Clicks, "Clicks should increment by 1")
}

func TestDirectLink_GetEvent_ValidEvent(t *testing.T) {
	// Arrange
	salonId := uuid.New()
	event := NewSalonInviteEvent(salonId)
	link := NewDirectLink(event)

	// Act
	retrievedEvent := link.Event

	// Assert
	assert.NotNil(t, retrievedEvent, "Retrieved event should not be nil")

	salonEvent, ok := retrievedEvent.(*SalonInviteEvent)
	assert.True(t, ok, "Retrieved event should be of type SalonInviteEvent")
	assert.Equal(t, salonId, salonEvent.SalonId, "SalonId should match")
}
