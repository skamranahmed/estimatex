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

// AskForVoteEventData represents data specific to the "ASK_FOR_VOTE" event
type AskForVoteEventData struct {
	TicketID string `json:"ticket_id"`
}

// MemberVotedEventData represents data specific to the "MEMBER_VOTED" event
type MemberVotedEventData struct {
	TicketID string `json:"ticket_id"`
	Vote     string `json:"vote"`
}

// VotingCompletedEventData represents data specific to the "VOTING_COMPLETED" event
type VotingCompletedEventData struct {
	TicketID string `json:"ticket_id"`
	Message  string `json:"message"`
}

// RevealVotesPromptEventData represents data specific to the "REVEAL_VOTES_PROMPT" event
type RevealVotesPromptEventData struct {
	TicketID string `json:"ticket_id"`
	Message  string `json:"message"`
}

// RevealVotesEventData represents data specific to the "REVEAL_VOTES" event
type RevealVotesEventData struct {
	TicketID string `json:"ticket_id"`
}
