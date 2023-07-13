package commands

import "github.com/maxguuse/gguuse-streams/internal/dataaccess"

type defaultCommand struct {
	cmds    dataaccess.CommandsRepository
	command string
}

func NewDefaultCommand(
	cmds dataaccess.CommandsRepository,
	command string,
) *defaultCommand {
	return &defaultCommand{
		cmds:    cmds,
		command: command,
	}
}

func (c *defaultCommand) GetAnswer() string {
	return c.cmds.GetCommandResponse(c.command)
}
