package resp

import (
	"encoding/json"
	"fmt"

	"redirectServer/internal/domain"
)

type DirectLinkDTO struct {
	NanoId  domain.NanoID     `json:"nanoId"`
	Payload map[string]string `json:"payload"`
	Event   string            `json:"event"`
}

func NewDirectLinkDTO(link domain.DirectLink) (*DirectLinkDTO, error) {
	var rawPayload map[string]string
	if err := json.Unmarshal(link.Payload, &rawPayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &DirectLinkDTO{
		NanoId:  link.NanoId,
		Payload: rawPayload,
		Event:   link.Event,
	}, nil
}
