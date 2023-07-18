package handlers

import (
	"log"
	"strings"
	"time"

	twitch_config "github.com/maxguuse/gguuse-streams/configs/twitch"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/maxguuse/gguuse-streams/internal/commands"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
	"golang.org/x/exp/slices"
)

type privateMessageHandler struct {
	cmds dataaccess.CommandsRepository
	anns dataaccess.AnnouncementsRepository
}

func NewPrivateMessageHandler(
	twitchCmds dataaccess.CommandsRepository,
	twitchAnns dataaccess.AnnouncementsRepository,
) *privateMessageHandler {
	return &privateMessageHandler{
		cmds: twitchCmds,
		anns: twitchAnns,
	}
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
		"help":             commands.NewHelpCommand(predefinedUserCommands, h.cmds.GetCommands()),
		"setmessage":       commands.NewSetMessageCommand(h.cmds, commandArgs),
		"newannouncement":  commands.NewNewAnnouncementCommand(h.anns, commandArgs),
		"stopannouncement": commands.NewStopAnnouncementCommand(h.anns, commandArgs),
		"title":            commands.NewSetTitleCommand(commandArgs),
	}

	commandHandler, ok := commandsHandlers[commandFromMessage]
	if !ok {
		commandHandler = commands.NewDefaultCommand(h.cmds, commandFromMessage)
	}
	answer := commandHandler.GetAnswer()

	if answer == "" {
		log.Printf("No such command: %s", commandFromMessage)
	} else {
		log.Printf("Replied with: %s", answer)
		twitch_config.IrcClient.Reply(twitch_config.Channel, m.ID, answer)
	}
}
