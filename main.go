package main

import (
	"fmt"
	"os"

	logger "log"

	"github.com/common-nighthawk/go-figure"
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

	choice := prompt.StringPrompt(
		"📚 Choose an option:\n\n    [1] Create a room\n    [2] Join a room\n\nType your choice (1 or 2):",
	)

	action, ok := responseActionMap[choice]
	if !ok {
		return ""
	}
	return action
}
