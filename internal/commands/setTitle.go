package commands

import (
	"errors"
	"fmt"
	"log"
	"strings"

	twitch_config "github.com/maxguuse/gguuse-streams/configs/twitch"
	"github.com/nicklaw5/helix/v2"
)

type newSetTitleCommand struct {
	cmdArgs []string
}

func NewSetTitleCommand(
	cmdArgs []string,
) *newSetTitleCommand {
	return &newSetTitleCommand{
		cmdArgs: cmdArgs,
	}
}

func (c *newSetTitleCommand) GetAnswer() (string, error) {
	if len(c.cmdArgs) < 1 {
		return "Usage: !title <title>", nil
	}

	usersResp, err := twitch_config.ApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{twitch_config.Channel},
	})

	if err != nil {
		log.Fatalf("Couldn't make GetUsers request to Twitch API, error: %s", err)
	} else if usersResp.ErrorMessage != "" {
		return "", errors.New(usersResp.ErrorMessage)
	}

	broadcasterId := usersResp.Data.Users[0].ID

	title := strings.Join(c.cmdArgs[0:], " ")

	editTitleResp, err := twitch_config.ApiClient.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: broadcasterId,
		Title:         title,
	})
	if err != nil {
		log.Fatalf("Couldn't make EditChannelInformation request to Twitch API, error: %s", err)
	} else if editTitleResp.ErrorMessage != "" {
		return "", errors.New(editTitleResp.ErrorMessage)
	}

	return fmt.Sprintf("Stream title changed to: %s", title), nil
}
