package commands

import (
	"fmt"
	"strings"

	"github.com/maxguuse/gguuse-streams/configs/repositories"
)

type setMessageCommand struct {
	cmdArgs []string
}

func NewSetMessageCommand(cmdArgs []string) *setMessageCommand {
	return &setMessageCommand{
		cmdArgs: cmdArgs,
	}
}

func (c *setMessageCommand) GetAnswer() string {
	if len(c.cmdArgs) < 1 {
		return "Usage: !setmessage <command> <message>"
	}

	cmdToChange := c.cmdArgs[0]
	newAnswer := strings.Join(c.cmdArgs[1:], " ")

	repositories.Commands.UpdateCommand(cmdToChange, newAnswer)
	repositories.Commands.SaveCommands()
	return fmt.Sprintf("Changed command: %s. New answer: %s", cmdToChange, newAnswer)
}
