package main

import (
	"fmt"
	"net/url"
	"os"

	logger "log"

	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/websocket"
	"github.com/skamranahmed/estimatex-client/prompt"
)

// set this to true for development, keep it false for production
const isDevelopment = true

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

	err := connectToEstimateXWebSocketEndpoint(action)
	if err != nil {
		log.Println(err.Error())
		return
	}
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

func connectToEstimateXWebSocketEndpoint(action UserAction) error {
	wsEndpoint := getWebSocketEndpoint()

	var err error
	var wsConnection *websocket.Conn

	if action == JOIN_ROOM {
		wsConnection, err = handleJoinRoomAction(wsEndpoint)
	} else if action == CREATE_ROOM {
		wsConnection, err = handleCreateRoomAction(wsEndpoint)
	}

	if err != nil {
		return err
	}

	defer closeWebSockeConnection(wsConnection)

	return nil
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
		Host:   "",    // TODO: to be filled when I will get the production URL of EstimateX server after deployment
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

	// TODO: need to inform the server that a room join event has occurred
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
	}
}
