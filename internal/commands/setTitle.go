package commands

import (
	"log"
	"strings"

	twitch_config "github.com/maxguuse/gguuse-streams/configs/twitch"
	"github.com/nicklaw5/helix/v2"
)

type newSetTitleCommand struct {
	cmdArgs []string
	channel string
}

func NewSetTitleCommand(
	cmdArgs []string,
	channel string,
) *newSetTitleCommand {
	return &newSetTitleCommand{
		cmdArgs: cmdArgs,
		channel: channel,
	}
}

func (c *newSetTitleCommand) GetAnswer() string {
	if len(c.cmdArgs) < 1 {
		return "Usage: !title <title>"
	}

	usersResp, err := twitch_config.ApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{c.channel},
	})
	if err != nil {
		log.Printf("Error fetching broadcaster id, error: %s", err)
		return ""
	}
	broadcasterId := usersResp.Data.Users[0].ID

	title := strings.Join(c.cmdArgs[0:], " ")

	_, err = twitch_config.ApiClient.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: broadcasterId,
		Title:         title,
	})
	if err != nil {
		log.Printf("Error changing stream title, error: %s", err)
		return ""
	}

	return ""
}
