package event

import "encoding/json"

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// RoomJoinEventData represents data specific to the "JOIN_ROOM" event
type RoomJoinEventData struct {
	RoomID string `json:"room_id"`
}
