package commands

import (
	"fmt"
	"strings"

	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
)

type setMessageCommand struct {
	cmds    dataaccess.CommandsRepository
	cmdArgs []string
}

func NewSetMessageCommand(
	cmds dataaccess.CommandsRepository,
	cmdArgs []string,
) *setMessageCommand {
	return &setMessageCommand{
		cmds:    cmds,
		cmdArgs: cmdArgs,
	}
}

func (c *setMessageCommand) GetAnswer() string {
	if len(c.cmdArgs) < 1 {
		return "Usage: !setmessage <command> <message>"
	}

	cmdToChange := c.cmdArgs[0]
	newCommand := strings.Join(c.cmdArgs[1:], " ")

	c.cmds.UpdateCommand(cmdToChange, newCommand)
	return fmt.Sprintf("Changed command: %s. New answer: %s", cmdToChange, newCommand)
}
