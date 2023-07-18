package handlers

import (
	"log"
	"strings"
	"time"

	"github.com/maxguuse/gguuse-streams/configs/repositories"
	twitch_config "github.com/maxguuse/gguuse-streams/configs/twitch"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/maxguuse/gguuse-streams/internal/commands"
	"golang.org/x/exp/slices"
)

type privateMessageHandler struct{}

func NewPrivateMessageHandler() *privateMessageHandler {
	return &privateMessageHandler{}
}

func (h *privateMessageHandler) Handle(m twitch.PrivateMessage) {
	if m.Message[0] != '!' {
		return
	}

	time.Sleep(time.Second)

	adminCommands := []string{
		"setmessage", "newannouncement", "stopannouncement", "title",
	}
	commandFromMessage := strings.Split(m.Message[1:], " ")[0]
	commandArgs := strings.Split(m.Message[1:], " ")[1:]

	log.Printf("Got command '%s' from %s", commandFromMessage, m.User.DisplayName)

	_, isBroadcaster := m.User.Badges["broadcaster"]
	_, isModerator := m.User.Badges["moderator"]

	hasAdminAccess := isBroadcaster || isModerator

	if !hasAdminAccess && slices.Contains(adminCommands, commandFromMessage) {
		log.Printf("%s doesn't have admin access", m.User.DisplayName)
		return
	}

	predefinedUserCommands := []string{
		"help",
	}

	commandsHandlers := map[string]commands.Command{
		"help":             commands.NewHelpCommand(predefinedUserCommands, repositories.Commands.GetCommands()),
		"setmessage":       commands.NewSetMessageCommand(commandArgs),
		"newannouncement":  commands.NewNewAnnouncementCommand(commandArgs),
		"stopannouncement": commands.NewStopAnnouncementCommand(commandArgs),
		"title":            commands.NewSetTitleCommand(commandArgs),
	}

	commandHandler, ok := commandsHandlers[commandFromMessage]
	if !ok {
		commandHandler = commands.NewDefaultCommand(commandFromMessage)
	}
	answer, err := commandHandler.GetAnswer()

	if answer != "" {
		log.Printf("Replied with: %s", answer)
		twitch_config.IrcClient.Reply(twitch_config.Channel, m.ID, answer)
	} else if err != nil {
		log.Printf("Couldn't answer to command, error occured: %s", err)
	} else {
		log.Println("Unknown error")
	}
}
