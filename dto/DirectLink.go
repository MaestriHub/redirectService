package dto

import "encoding/json"

type DirectLink struct { //link.maetry.com/{NanoId
	ID      string          // NanoId
	Payload json.RawMessage `json:"payload" gorm:"type:jsonb"`
	Event   string          `json:"event"`
}
