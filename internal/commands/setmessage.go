package commands

import (
	"fmt"
	"log"
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

func (c *setMessageCommand) GetAnswer() (string, error) {
	if len(c.cmdArgs) < 1 {
		return "Usage: !setmessage <command> <message>", nil
	}

	cmdToChange := c.cmdArgs[0]
	newAnswer := strings.Join(c.cmdArgs[1:], " ")

	repositories.Commands.UpdateCommand(cmdToChange, newAnswer)
	err := repositories.Commands.SaveCommands()
	if err != nil {
		log.Fatalf("Error occured while saving commands: %s", err)
	}

	return fmt.Sprintf("Changed command: %s. New answer: %s", cmdToChange, newAnswer), nil
}
