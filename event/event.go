package event

import "encoding/json"

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// CreateRoomEventData represents data specific to the "CREATE_ROOM" event
type CreateRoomEventData struct {
	RoomID string `json:"room_id"`
}

// RoomJoinEventData represents data specific to the "JOIN_ROOM" event
type RoomJoinEventData struct {
	RoomID string `json:"room_id"`
}

// RoomJoinUpdatesEventData represents data specific to the "ROOM_JOIN_UPDATES" event
type RoomJoinUpdatesEventData struct {
	Message string `json:"message"`
}

// RoomCapacityReachedEventData represents data specific to the "ROOM_CAPACITY_REACHED" event
type RoomCapacityReachedEventData struct {
	Message string `json:"message"`
}

// BeginVotingPromptEventData represents data specific to the "BEGIN_VOTING_PROMPT" event
type BeginVotingPromptEventData struct {
	Message string `json:"message"`
}

// BeginVotingEventData represents data specific to the "BEGIN_VOTING" event
type BeginVotingEventData struct {
	TicketID string `json:"ticket_id"`
}
