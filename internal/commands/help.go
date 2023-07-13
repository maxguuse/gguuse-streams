package commands

import (
	"fmt"
	"strings"
)

type helpCommand struct {
	cmds []string
}

func NewHelpCommand(commandsSources ...[]string) *helpCommand {
	cmdsFromSources := []string{}
	for i := 0; i < len(commandsSources); i++ {
		cmdsFromSources = append(cmdsFromSources, commandsSources[i]...)
	}

	return &helpCommand{
		cmds: cmdsFromSources,
	}
}

func (c *helpCommand) GetAnswer() string {
	return fmt.Sprintf("Available commands: %s", strings.Join(c.cmds, " "))
}
