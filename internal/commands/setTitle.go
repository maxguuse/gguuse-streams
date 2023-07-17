package commands

import (
	"strings"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/maxguuse/gguuse-streams/tools"

)

type newSetTitleCommand struct {
	cmdArgs []string
	client  *twitch.Client
	channel string
}

func NewSetTitleCommand(
	cmdArgs []string,
	client *twitch.Client,
	channel string,
) *newSetTitleCommand {
	return &newSetTitleCommand{
		cmdArgs: cmdArgs,
		client:  client,
		channel: channel,
	}
}

func (c *newSetTitleCommand) GetAnswer() string{
	if (len(c.cmdArgs) < 1){
		return "Usage: !title <title>"
	}
	broadcasterId := tools.GetUserId(c.channel)
	title := strings.Join(c.cmdArgs[0:], " ")
	tools.SetTitle(broadcasterId, title)
	return ""
}