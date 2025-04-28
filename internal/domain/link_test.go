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
	link, err := NewDirectLink(event)
	if err != nil {
		t.Fatalf("Failed to create DirectLink: %v", err)
	}

	// Assert
	assert.NotNil(t, link)
	assert.NotEmpty(t, link.NanoId, "NanoId should not be empty")
	assert.Equal(t, 0, link.Clicks, "Clicks should be initialized to 0")
	assert.Equal(t, event.GetType(), link.Event, "Event type should match")

	// Check Payload deserialization
	var payloadEvent SalonInviteEvent
	err = json.Unmarshal(link.Payload, &payloadEvent)
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

func TestDirectLink_ToURL(t *testing.T) {
	// Arrange
	link := DirectLink{NanoId: "abc123"}

	// Act
	url := link.ToURL()

	// Assert
	expectedURL := "https://link.maetry.com/abc123"
	assert.Equal(t, expectedURL, url, "URL should match the expected format")
}

func TestDirectLink_GetEvent_ValidEvent(t *testing.T) {
	// Arrange
	salonId := uuid.New()
	event := NewSalonInviteEvent(salonId)
	link, err := NewDirectLink(event)
	if err != nil {
		t.Fatalf("Failed to create DirectLink: %v", err)
	}

	// Act
	retrievedEvent, err := link.GetEvent()

	// Assert
	assert.NoError(t, err, "GetEvent should not return an error for valid events")
	assert.NotNil(t, retrievedEvent, "Retrieved event should not be nil")

	salonEvent, ok := retrievedEvent.(*SalonInviteEvent)
	assert.True(t, ok, "Retrieved event should be of type SalonInviteEvent")
	assert.Equal(t, salonId, salonEvent.SalonId, "SalonId should match")
}
