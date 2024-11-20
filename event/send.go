package event

import (
	"encoding/json"

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

func SendBeginVotingEvent(wsConnection *websocket.Conn, ticketId string) {
	beginVotingEventData := BeginVotingEventData{
		TicketID: ticketId,
	}

	beginVotingEventJsonData, _ := json.Marshal(beginVotingEventData)
	beginVotingEvent := Event{
		Type: "BEGIN_VOTING",
		Data: json.RawMessage(beginVotingEventJsonData),
	}

	sendMessage(wsConnection, beginVotingEvent)
}

func SendMemberVotedEvent(wsConnection *websocket.Conn, ticketId string, vote string) {
	memberVotedEventData := MemberVotedEventData{
		TicketID: ticketId,
		Vote:     vote,
	}

	memberVotedEventJsonData, _ := json.Marshal(memberVotedEventData)
	beginVotingEvent := Event{
		Type: "MEMBER_VOTED",
		Data: json.RawMessage(memberVotedEventJsonData),
	}

	sendMessage(wsConnection, beginVotingEvent)
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
