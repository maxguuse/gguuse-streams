package handlers

import (
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/maxguuse/gguuse-streams/internal/commands"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
	"golang.org/x/exp/slices"
)

type privateMessageHandler struct {
	client  *twitch.Client
	channel string

	cmds dataaccess.CommandsRepository
	anns dataaccess.AnnouncmentsRepository
}

func NewPrivateMessageHandler(
	twitchClient *twitch.Client,
	twitchChannel string,
	twitchCmds dataaccess.CommandsRepository,
	twitchAnns dataaccess.AnnouncmentsRepository,
) *privateMessageHandler {
	return &privateMessageHandler{
		client:  twitchClient,
		channel: twitchChannel,
		cmds:    twitchCmds,
		anns:    twitchAnns,
	}
}

func (h *privateMessageHandler) Handle(m twitch.PrivateMessage) {
	if m.Message[0] != '!' {
		return
	}

	time.Sleep(time.Second)

	adminCommands := []string{
		"setmessage", "newannouncment", "stopannouncment",
	}
	commandFromMessage := strings.Split(m.Message[1:], " ")[0]
	commandArgs := strings.Split(m.Message[1:], " ")[1:]

	_, isBroadcaster := m.User.Badges["broadcaster"]
	_, isModerator := m.User.Badges["moderator"]

	hasAdminAccess := isBroadcaster || isModerator

	if !hasAdminAccess && slices.Contains(adminCommands, commandFromMessage) {
		return
	}

	predefinedUserCommands := []string{
		"help",
	}

	commandsHandlers := map[string]commands.Command{
		"help":            commands.NewHelpCommand(predefinedUserCommands, h.cmds.GetCommands()),
		"setmessage":      commands.NewSetMessageCommand(h.cmds, commandArgs),
		"newannouncment":  commands.NewNewAnnouncmentCommand(h.anns, commandArgs, h.client, h.channel),
		"stopannouncment": commands.NewStopAnnouncmentCommand(h.anns, commandArgs),
	}

	commandHandler, ok := commandsHandlers[commandFromMessage]
	if !ok {
		commandHandler = commands.NewDefaultCommand(h.cmds, commandFromMessage)
	}
	answer := commandHandler.GetAnswer()

	h.client.Reply(h.channel, m.ID, answer)
}
