package handlers

import (
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/maxguuse/gguuse-streams/internal/commands"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
	"golang.org/x/exp/slices"
)

type privateMessageHandler struct {
	client  *twitch.Client
	channel string

	cmds dataaccess.CommandsRepository
}

func NewPrivateMessageHandler(
	twitchClient *twitch.Client,
	twitchChannel string,
	twitchCmds dataaccess.CommandsRepository,
) *privateMessageHandler {
	return &privateMessageHandler{
		client:  twitchClient,
		channel: twitchChannel,
		cmds:    twitchCmds,
	}
}

func (h *privateMessageHandler) Handle(m twitch.PrivateMessage) {
	if m.Message[0] != '!' {
		return
	}

	adminCommands := []string{
		"setmessage",
	}
	commandFromMessage := strings.Split(m.Message[1:], " ")[0]
	if m.User.DisplayName != "GGuuse" && slices.Contains(adminCommands, commandFromMessage) {
		return
	}

	predefinedUserCommands := []string{
		"help",
	}

	commandsHandlers := map[string]commands.Command{
		"help":       commands.NewHelpCommand(predefinedUserCommands, h.cmds.GetCommands()),
		"setmessage": commands.NewSetMessageCommand(h.cmds, strings.Split(m.Message[1:], " ")[1:]),
	}

	commandHandler, ok := commandsHandlers[commandFromMessage]
	if !ok {
		commandHandler = commands.NewDefaultCommand(h.cmds, commandFromMessage)
	}
	answer := commandHandler.GetAnswer()

	h.client.Reply(h.channel, m.ID, answer)
}
