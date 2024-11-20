package event

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	logger "log"

	"github.com/gorilla/websocket"
	"github.com/olekukonko/tablewriter"
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
	eventHandlers["ASK_FOR_VOTE"] = AskForVoteEventHandler
	eventHandlers["VOTING_COMPLETED"] = VotingCompletedEventHandler
	eventHandlers["REVEAL_VOTES_PROMPT"] = RevealVotesPromptEventHandler
	eventHandlers["VOTES_REVEALED"] = VotesRevealedEventHandler
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

func AskForVoteEventHandler(wsConnection *websocket.Conn, event Event) error {
	var askForVoteEventData AskForVoteEventData
	err := json.Unmarshal(event.Data, &askForVoteEventData)
	if err != nil {
		log.Println("unable to handle ASK_FOR_VOTE event", err)
		return nil
	}

	askForVoteMessage := fmt.Sprintf(
		"> 📝 Choose a story point for the ticket: %s\n\n    [1] 1\n    [2] 2\n    [3] 3\n    [5] 5\n    [8] 8\n    [13] 13\n    [21] 21\n\nType your choice (1 or 2 or 3 or 5 or 8 or 13 or 21):", askForVoteEventData.TicketID,
	)
	vote := prompt.StringInputPrompt(askForVoteMessage)

	log.Printf("👍 You voted %v for the ticket id: %s\n", vote, askForVoteEventData.TicketID)

	SendMemberVotedEvent(wsConnection, askForVoteEventData.TicketID, vote)
	return nil
}

func VotingCompletedEventHandler(wsConnection *websocket.Conn, event Event) error {
	var votingCompletedEventData VotingCompletedEventData
	err := json.Unmarshal(event.Data, &votingCompletedEventData)
	if err != nil {
		log.Println("unable to handle VOTING_COMPLETED event", err)
		return nil
	}

	// the "VOTING_COMPLETED" event message from the server will be plain text,
	// hence, we will simply log it and use it as the user prompt text
	log.Println(votingCompletedEventData.Message)
	return nil
}

func RevealVotesPromptEventHandler(wsConnection *websocket.Conn, event Event) error {
	var revealVotesPromptEventData RevealVotesPromptEventData
	err := json.Unmarshal(event.Data, &revealVotesPromptEventData)
	if err != nil {
		log.Println("unable to handle REVEAL_VOTES_PROMPT event", err)
		return nil
	}

	choice := prompt.StringInputPrompt(
		"✨ Time to reveal the votes.\n\nEnter 'Y' to confirm:",
	)
	if choice != "Y" {
		return fmt.Errorf("you chose not to reveal the votes")
	}

	SendRevealVotesEvent(wsConnection, revealVotesPromptEventData.TicketID)
	return nil
}

func VotesRevealedEventHandler(serverWsConnection *websocket.Conn, event Event) error {
	var votesRevealedEventData VotesRevealedEventData
	err := json.Unmarshal(event.Data, &votesRevealedEventData)
	if err != nil {
		log.Println("unable to handle VOTES_REVEALED event", err)
		return nil
	}

	// the "VOTES_REVEALED" event message from the server will be plain text,
	// hence, we will simply log it and use it as the user prompt text
	log.Printf("👇 Votes for the ticket id: %v", votesRevealedEventData.TicketID)

	renderVotes(votesRevealedEventData.ClientVoteChoiceMap)
	return nil
}

func renderVotes(voteChoiceMap map[string]Vote) {
	grouped := make(map[string][]Vote)

	for _, vote := range voteChoiceMap {
		grouped[vote.Value] = append(grouped[vote.Value], vote)
	}

	var voteValues []string
	for voteValue := range grouped {
		voteValues = append(voteValues, voteValue)
	}
	sort.Strings(voteValues)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Vote Value", "Member(s)", "Count"})

	for _, voteValue := range voteValues {
		votes := grouped[voteValue]
		voteCount := fmt.Sprintf("%v", len(votes))
		for index, vote := range votes {
			clientName := voteChoiceMap[vote.MemberID].MemberName
			if index == 0 {
				table.Append([]string{voteValue, clientName, voteCount})
			} else {
				table.Append([]string{"", clientName, ""})
			}
		}
		table.Append([]string{"", ""})
	}

	table.Render()
}
