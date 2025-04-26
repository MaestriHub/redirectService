package domain

import (
	"encoding/json"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type NanoID = string
type Payload = []byte

type DirectLink struct {
	NanoId  NanoID
	Clicks  int     `gorm:"default:0"`
	Payload Payload `gorm:"type:jsonb"`
	Event   string
}

func NewDirectLink(event InviteEvent) (*DirectLink, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("error marshaling to JSON: %v", err)
	}

	return &DirectLink{
		NanoId:  gonanoid.Must(8),
		Clicks:  0,
		Payload: payload,
		Event:   event.GetType(),
	}, nil
}

func (l *DirectLink) IncClicks() {
	l.Clicks++
}

func (l *DirectLink) ToURL() string {
	return "https://link.maetry.com/" + l.NanoId
}

func (l *DirectLink) Validate() error {
	return nil
}

func (l *DirectLink) GetEvent() (InviteEvent, error) {
	switch l.Event {
	case EmployeeInvite:
		var event EmployeeInviteEvent
		if err := json.Unmarshal(l.Payload, &event); err != nil {
			return nil, err
		}
		event.BaseInviteEvent = BaseInviteEvent{Type: l.Event}
		return &event, nil
	case SalonInvite:
		var event SalonInviteEvent
		if err := json.Unmarshal(l.Payload, &event); err != nil {
			return nil, err
		}
		event.BaseInviteEvent = BaseInviteEvent{Type: l.Event}
		return &event, nil
	case ClientInvite:
		var event ClientInviteEvent
		if err := json.Unmarshal(l.Payload, &event); err != nil {
			return nil, err
		}
		event.BaseInviteEvent = BaseInviteEvent{Type: l.Event}
		return &event, nil
	default:
		return nil, fmt.Errorf("unknown event: %s; id - %s", l.Event, l.NanoId)
	}
}
