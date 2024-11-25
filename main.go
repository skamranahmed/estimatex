package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"

	logger "log"

	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/websocket"
	"github.com/skamranahmed/estimatex/event"
	"github.com/skamranahmed/estimatex/prompt"
)

// set this to true for development, keep it false for production
const isDevelopment = false

var (
	log = logger.New(os.Stdout, "> ", 0) // 0 means no flags (no time, no date, etc.)
)

type UserAction string

const (
	CREATE_ROOM UserAction = "CREATE_ROOM"
	JOIN_ROOM   UserAction = "JOIN_ROOM"
)

func main() {
	displayWelcomeMessage()

	action := promptUserAction()
	if action == "" {
		log.Println("❌ Exiting program due to invalid choice.")
		return
	}

	event.SetupEventHandlers()

	wsConnection, err := connectToEstimateXWebSocketEndpoint(action)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer closeWebSockeConnection(wsConnection)

	// read messages from the server, spawning a go-routine for this, because it is blocking in nature
	go readMessages(wsConnection)

	// wait for interrupt signal from the user to exit the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	log.Println("👋 Exiting...")
}

func displayWelcomeMessage() {
	// generate an ASCII art
	fig := figure.NewFigure("EstimateX", "", true)
	fig.Print()

	if isDevelopment {
		fmt.Println("\n🛠️🚀 Development Mode")
	}

	fmt.Println("\n👋 Welcome! Use the tool to proceed. Press CTRL + C to exit when you're done.\n")
}

func promptUserAction() UserAction {
	responseActionMap := map[string]UserAction{
		"1": CREATE_ROOM,
		"2": JOIN_ROOM,
	}

	choice := prompt.StringInputPrompt(
		"📚 Choose an option:\n\n    [1] Create a room\n    [2] Join a room\n\nType your choice (1 or 2):",
	)

	action, ok := responseActionMap[choice]
	if !ok {
		return ""
	}
	return action
}

func connectToEstimateXWebSocketEndpoint(action UserAction) (*websocket.Conn, error) {
	wsEndpoint := getWebSocketEndpoint()

	var err error
	var wsConnection *websocket.Conn

	if action == JOIN_ROOM {
		wsConnection, err = handleJoinRoomAction(wsEndpoint)
	} else if action == CREATE_ROOM {
		wsConnection, err = handleCreateRoomAction(wsEndpoint)
	}

	if err != nil {
		return nil, err
	}

	return wsConnection, nil
}

func readMessages(wsConnection *websocket.Conn) {
	defer closeWebSockeConnection(wsConnection)

	// continuously read messages from the server
	for {
		_, payload, err := wsConnection.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				// server sent the "Close Handshake" message
				log.Println("🔌 Server closed the connection. Exiting program...")
				return
			}

			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				// server closed the WebSocket connection abruptly
				log.Println("💔 Server closed the connection. Exiting program...")
				return
			}

			log.Println("❌ Server closed the connection. Exiting program...")
			return
		}

		// process the received message(s)
		var receivedEvent event.Event
		err = json.Unmarshal(payload, &receivedEvent)
		if err != nil {
			// this means the received message must be in plain text, so just log it
			log.Println(string(payload))
			continue
		}

		// handle received event message
		err = event.HandleEvent(wsConnection, receivedEvent)
		if err != nil {
			log.Printf("❗️ Error while handling the received event %v: %+v Exiting program...", receivedEvent.Type, err)
			return
		}
	}
}

func getWebSocketEndpoint() url.URL {
	path := "/ws"
	if isDevelopment {
		return url.URL{
			Scheme: "ws",
			Host:   "localhost:8080",
			Path:   path,
		}
	}

	return url.URL{
		Scheme: "wss", // use websocket secure protocol in production
		Host:   "estimatex-server-production.up.railway.app",
		Path:   path,
	}
}

func handleJoinRoomAction(wsEndpoint url.URL) (*websocket.Conn, error) {
	roomId := prompt.StringInputPrompt("> 📝 Enter the room id which you would like to join:")

	clientName := prompt.StringInputPrompt("> 📝 Enter your name:")
	if clientName == "" {
		return nil, fmt.Errorf("❌ Empty client name is not allowed. Exiting program...")
	}

	wsEndpoint.RawQuery = fmt.Sprintf("action=%s&name=%s&room_id=%s", JOIN_ROOM, url.QueryEscape(clientName), url.QueryEscape(roomId))

	wsConnection, err := establishWebSocketConnection(wsEndpoint)
	if err != nil {
		return nil, err
	}

	// inform the server that a room join event has occurred
	event.SendRoomJoinEvent(wsConnection, roomId)
	return wsConnection, nil
}

func handleCreateRoomAction(wsEndpoint url.URL) (*websocket.Conn, error) {
	maxRoomCapacity, err := prompt.IntegerInputPrompt("> 📝 Enter the room max capacity:")
	if err != nil {
		return nil, fmt.Errorf("❌ You have entered an invalid numerical value for room max capacity. Exiting program...")
	}

	clientName := prompt.StringInputPrompt("> 📝 Enter your name:")
	if clientName == "" {
		return nil, fmt.Errorf("❌ Empty client name is not allowed. Exiting program...")
	}

	wsEndpoint.RawQuery = fmt.Sprintf("action=%s&name=%s&max_room_capacity=%d", CREATE_ROOM, url.QueryEscape(clientName), maxRoomCapacity)

	wsConnection, err := establishWebSocketConnection(wsEndpoint)
	if err != nil {
		return nil, err
	}

	return wsConnection, nil
}

func establishWebSocketConnection(wsEndpoint url.URL) (*websocket.Conn, error) {
	wsConnection, _, err := websocket.DefaultDialer.Dial(wsEndpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("❌ Got an error while trying to connect to EstimateX endpoint: %v", err)
	}
	return wsConnection, nil
}

func closeWebSockeConnection(wsConnection *websocket.Conn) {
	if wsConnection != nil {
		// gracefully close the WebSocket connection by sending a "Close Handshake" message from the client
		wsConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Client closing connection"))
		wsConnection.Close()
		os.Exit(0)
	}
}
