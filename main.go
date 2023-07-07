package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
	"github.com/sunspots/tmi"
)

var (
	connection *tmi.Connection
	commands   dataaccess.CommandsRepository
)

func receiveMessages() {
	err := error(nil)
	for err == nil {
		err = readMessage(connection)
	}
}

func readMessage(connection *tmi.Connection) error {
	msg, err := connection.ReadMessage()
	if err != nil {
		log.Println(err)
		return err
	}
	if msg.Command == "PRIVMSG" {
		processMessage(msg)
	}
	return nil
}

func processMessage(msg *tmi.Message) {
	var answer string

	if msg.Trailing[0] == '!' {
		msgText := msg.Trailing
		msgTags := msg.Tags
		msgChannel := msg.Channel()

		answer = processCommand(msgText[1:], msgTags, msgChannel)
	}

	if answer != "" {
		time.Sleep(1 * time.Second)
		connection.Send(answer)
	}
}

func processCommand(commandMsg string, tags map[string]string, channel string) (answer string) {
	commandMsgSlice := strings.Split(commandMsg, " ")

	command := commandMsgSlice[0]
	//args := commandMsgSlice[1:]

	answerWith := func(msg string) string {
		return fmt.Sprintf("@reply-parent-msg-id=%s PRIVMSG %s :%s", tags["id"], channel, msg)
	}

	switch command {
	case "setmessage":
		if tags["display-name"] != "GGuuse" {
			return ""
		}
		commandToChange := commandMsgSlice[1]
		newMessage := strings.Join(commandMsgSlice[2:], " ")

		commands.UpdateCommand(commandToChange, newMessage)
		commands.SaveCommands()
		return answerWith(fmt.Sprintf("%s: %s", commandToChange, newMessage))
	case "help":
		cmds := commands.GetCommands()

		return answerWith(fmt.Sprintf("Available commands: %s", strings.Join(cmds, " ")))
	default:
		answer := commands.GetCommandResponse(command)
		return answerWith(answer)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()

	username := os.Getenv("USERNAME")
	authToken := os.Getenv("AUTH_TOKEN")
	connection = tmi.Connect(username, authToken)

	channelName := os.Args[1]
	channel := fmt.Sprintf("#%s", channelName)
	connection.Join(channel)

	// TODO: Replace with logs
	connection.Debug = true

	commands = dataaccess.NewCommandsRepository()
	commands.LoadCommands()

	receiveMessages()
}
