package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"redirectServer/internal/domain"
)

func TestDirectLink_ToURL(t *testing.T) {
	// Arrange
	urlPrefix := "http://хуй.до.ушей/"
	nanoId := domain.NanoID("abc123")
	service := NewUrlGeneratorService(urlPrefix)
	link := domain.DirectLink{NanoId: nanoId}

	// Act
	url := service.Generate(&link)

	// Assert
	expectedURL := urlPrefix + nanoId
	assert.Equal(t, expectedURL, url, "URL should match the expected format")
}
