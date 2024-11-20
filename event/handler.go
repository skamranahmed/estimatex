package event

import (
	"encoding/json"
	"fmt"
	"os"

	logger "log"

	"github.com/gorilla/websocket"
	"github.com/skamranahmed/estimatex-client/prompt"
)

type EventHandler func(wsConnection *websocket.Conn, event Event) error

var (
	log = logger.New(os.Stdout, "> ", 0) // 0 means no flags (no time, no date, etc.)

	eventHandlers = make(map[string]EventHandler)
)

func EventNotSupportedError(eventType string) error {
	return fmt.Errorf("%s event type is not supported yet.", eventType)
}

// sets up event handlers to process various types of event messages received from the server
func SetupEventHandlers() {
	eventHandlers["CREATE_ROOM"] = CreateRoomEventHandler
	eventHandlers["ROOM_JOIN_UPDATES"] = RoomJoinUpdatesEventHandler
	eventHandlers["ROOM_CAPACITY_REACHED"] = RoomCapacityReachedEventHandler
	eventHandlers["BEGIN_VOTING_PROMPT"] = BeginVotingPromptEventHandler
}

func HandleEvent(wsConnection *websocket.Conn, event Event) error {
	eventHandler, ok := eventHandlers[event.Type]
	if ok {
		err := eventHandler(wsConnection, event)
		if err != nil {
			return err
		}
		return nil
	} else {
		return EventNotSupportedError(event.Type)
	}
}

func CreateRoomEventHandler(wsConnection *websocket.Conn, event Event) error {
	var roomCreationEventData CreateRoomEventData
	err := json.Unmarshal(event.Data, &roomCreationEventData)
	if err != nil {
		log.Println("unable to handle CREATE_ROOM event", err)
		return nil
	}

	// upon receiving the "CREATE_ROOM" event message from the server,
	// we need to send the "JOIN_ROOM" event message to the server
	SendRoomJoinEvent(wsConnection, roomCreationEventData.RoomID)

	return nil
}

func RoomJoinUpdatesEventHandler(wsConnection *websocket.Conn, event Event) error {
	var roomJoinUpdatesEventData RoomJoinUpdatesEventData
	err := json.Unmarshal(event.Data, &roomJoinUpdatesEventData)
	if err != nil {
		log.Println("unable to handle ROOM_JOIN_UPDATES event", err)
		return nil
	}

	// the "ROOM_JOIN_UPDATES" event message from the server will be plain text,
	// hence, we will simply log it
	log.Println(roomJoinUpdatesEventData.Message)
	return nil
}

func RoomCapacityReachedEventHandler(wsConnection *websocket.Conn, event Event) error {
	var roomCapacityReachedEventData RoomCapacityReachedEventData
	err := json.Unmarshal(event.Data, &roomCapacityReachedEventData)
	if err != nil {
		log.Println("unable to handle ROOM_CAPACITY_REACHED event", err)
		return nil
	}

	// the "ROOM_CAPACITY_REACHED" event message from the server will be plain text,
	// hence, we will simply log it
	log.Println(roomCapacityReachedEventData.Message)
	return nil
}

func BeginVotingPromptEventHandler(wsConnection *websocket.Conn, event Event) error {
	var beginVotingPromptEventData BeginVotingPromptEventData
	err := json.Unmarshal(event.Data, &beginVotingPromptEventData)
	if err != nil {
		log.Println("unable to handle BEGIN_VOTING_PROMPT event", err)
		return nil
	}

	// the "BEGIN_VOTING_PROMPT" event message from the server will be plain text,
	// hence, we will simply log it and use it as the user prompt text
	messageToDisplay := fmt.Sprintf("> %s", beginVotingPromptEventData.Message)
	ticketId := prompt.StringInputPrompt(messageToDisplay)

	if ticketId == "" {
		return fmt.Errorf("empty ticket id is not allowed.")
	}

	SendBeginVotingEvent(wsConnection, ticketId)
	return nil
}
