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

type twitch_tags = map[string]string

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

func processCommand(commandMsg string, tags twitch_tags, channel string) (answer string) {
	command, args := parseCommand(commandMsg)

	replyWith := func(msg string) string {
		return fmt.Sprintf("@reply-parent-msg-id=%s PRIVMSG %s :%s", tags["id"], channel, msg)
	}

	handlers := map[string]func([]string, twitch_tags) string{
		"setmessage": handleSetMessage,
		"help":       handleHelp,
	}
	defaultHandler := func([]string, twitch_tags) string {
		return commands.GetCommandResponse(command)
	}

	handler, ok := handlers[command]
	if !ok {
		handler = defaultHandler
	}

	return replyWith(handler(args, tags))
}

func parseCommand(commandMsg string) (command string, args []string) {
	commandMsgSlice := strings.Split(commandMsg, " ")
	return commandMsgSlice[0], commandMsgSlice[1:]
}

func handleSetMessage(args []string, tags twitch_tags) string {
	if len(args) < 2 {
		return "Usage: !setmessage <command> <message>"
	}
	if tags["display-name"] != "GGuuse" {
		return ""
	}
	commandToChange := args[0]
	newMessage := strings.Join(args[1:], " ")

	commands.UpdateCommand(commandToChange, newMessage)
	return fmt.Sprintf("%s: %s", commandToChange, newMessage)
}

func handleHelp(args []string, tags twitch_tags) string {
	cmds := commands.GetCommands()

	return fmt.Sprintf("Available commands: %s", strings.Join(cmds, " "))
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()

	username := os.Getenv("NICK")
	authToken := os.Getenv("AUTH_TOKEN")
	connection = tmi.Connect(username, authToken)

	channelName := os.Args[1]
	channel := fmt.Sprintf("#%s", channelName)
	connection.Join(channel)

	// TODO: Replace with logs
	connection.Debug = true

	commands = dataaccess.NewCommandsRepository()

	err := commands.LoadCommands()
	if err == nil {
		receiveMessages()
	}
	commands.SaveCommands()
}
