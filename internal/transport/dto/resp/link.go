package resp

import (
	"encoding/json"
	"net/http"

	"redirectServer/internal/domain"
	"redirectServer/pkg"
)

type DirectLinkDTO struct {
	NanoId  domain.NanoID     `json:"nanoId"`
	Payload map[string]string `json:"payload"`
	Event   string            `json:"event"`
}

func NewDirectLinkDTO(link domain.DirectLink) (*DirectLinkDTO, *pkg.ErrorS) {
	payload, err := link.Event.GetPayload()
	if err != nil {
		return nil, pkg.NewErrorS("Oops", http.StatusInternalServerError)
	}

	var rawPayload map[string]string
	if err := json.Unmarshal(payload, &rawPayload); err != nil {
		return nil, pkg.NewErrorS("Ooops", http.StatusInternalServerError)
	}

	return &DirectLinkDTO{
		NanoId:  link.NanoId,
		Payload: rawPayload,
		Event:   link.Event.GetType(),
	}, nil
}
