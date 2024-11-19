package main

import (
	"fmt"
	"os"

	logger "log"

	"github.com/common-nighthawk/go-figure"
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
