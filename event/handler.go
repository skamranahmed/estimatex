package event

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type EventHandler func(wsConnection *websocket.Conn, event Event) error

var (
	eventHandlers = make(map[string]EventHandler)
)

func EventNotSupportedError(eventType string) error {
	return fmt.Errorf("%s event type is not supported yet.", eventType)
}

// sets up event handlers to process various types of event messages received from the server
func SetupEventHandlers() {
	eventHandlers["CREATE_ROOM"] = CreateRoomEventHandler
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
