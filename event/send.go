package event

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func SendRoomJoinEvent(wsConnection *websocket.Conn, roomId string) {
	roomJoiningEventData := RoomJoinEventData{
		RoomID: roomId,
	}

	roomJoiningEventJsonData, _ := json.Marshal(roomJoiningEventData)
	roomJoiningEvent := Event{
		Type: "JOIN_ROOM",
		Data: json.RawMessage(roomJoiningEventJsonData),
	}

	sendMessage(wsConnection, roomJoiningEvent)
}

func sendMessage(wConnection *websocket.Conn, event Event) {
	jsonMessage, err := json.Marshal(event)
	if err != nil {
		log.Printf("unable to marshal message: %+v, error: %+v", event, err)
	}

	// send message to the server
	err = wConnection.WriteMessage(websocket.TextMessage, []byte(jsonMessage))
	if err != nil {
		log.Println(err)
		return
	}
}
