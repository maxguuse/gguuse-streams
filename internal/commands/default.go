package commands

import "github.com/maxguuse/gguuse-streams/configs/repositories"

type defaultCommand struct {
	command string
}

func NewDefaultCommand(
	command string,
) *defaultCommand {
	return &defaultCommand{
		command: command,
	}
}

func (c *defaultCommand) GetAnswer() string {
	return repositories.Commands.GetResponse(c.command)
}
