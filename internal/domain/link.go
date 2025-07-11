package domain

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type NanoID = string

type DirectLink struct {
	NanoId NanoID
	Clicks int
	Event  Event
}

func NewDirectLink(event Event) *DirectLink {
	return &DirectLink{
		NanoId: gonanoid.Must(8),
		Clicks: 0,
		Event:  event,
	}
}

func (l *DirectLink) IncClicks() {
	l.Clicks++
}

func (l *DirectLink) Validate() error {
	return nil
}
