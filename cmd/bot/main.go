package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sunspots/tmi"
)

func readLoop(connection *tmi.Connection) {
	err := error(nil)
	for err == nil {
		err = readMessage(connection)
	}
}

func readMessage(connection *tmi.Connection) error {
	evt, err := connection.ReadMessage()
	if err != nil {
		log.Println(err)
		return err
	}
	if evt.Command == "PRIVMSG" {
		processMessage(connection, evt)
	}
	return nil
}

func processMessage(connection *tmi.Connection, evt *tmi.Message) {
	var answer string

	if evt.Trailing[0] == '!' {
		msgText := evt.Trailing
		msgTags := evt.Tags
		msgChannel := evt.Channel()

		answer = processCommand(msgText[1:], msgTags, msgChannel)
	}

	if answer != "" {
		time.Sleep(1 * time.Second)
		connection.Send(answer)
	}
}

var commands map[string]string = map[string]string{
	"ping": "pong",
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

		commands[commandToChange] = newMessage
		return answerWith(fmt.Sprintf("%s: %s", commandToChange, newMessage))
	case "help":
		keys := make([]string, 0, len(commands))
		for k := range commands {
			keys = append(keys, k)
		}

		return answerWith(fmt.Sprintf("Available commands: %s", strings.Join(keys, " ")))
	default:
		answer, isExist := commands[command]
		if isExist {
			return answerWith(answer)
		}
		return ""
	}
}

func main() {
	var channel string
	fmt.Print("Введите название канала: ")
	fmt.Scanln(&channel)

	channel = fmt.Sprintf("#%s", channel)

	connection := tmi.Connect("gguuse", os.Getenv("GGUUSE_STREAMS_AUTH_TOKEN"))

	connection.Debug = true
	connection.Join(channel)

	readLoop(connection)
}
