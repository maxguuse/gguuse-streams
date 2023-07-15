package commands

import (
	"fmt"

	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
)

type stopAnnouncmentCommand struct {
	anns    dataaccess.AnnouncmentsRepository
	cmdArgs []string
}

func NewStopAnnouncmentCommand(
	anns dataaccess.AnnouncmentsRepository,
	cmdArgs []string,
) *stopAnnouncmentCommand {
	return &stopAnnouncmentCommand{
		anns:    anns,
		cmdArgs: cmdArgs,
	}
}

func (c *stopAnnouncmentCommand) GetAnswer() string {
	if len(c.cmdArgs) != 1 {
		return "Usage: !stopannouncment <id>"
	}

	ann, isExists := c.anns.GetAnnouncment(c.cmdArgs[0])
	if !isExists {
		return fmt.Sprintf("Announcment with ID: %s doesn't exist", c.cmdArgs[0])
	}

	c.anns.RemoveAnnouncment(ann.GetID())
	c.anns.SaveAnnouncments()

	return fmt.Sprintf("Announcment with ID: %s has been deleted", c.cmdArgs[0])
}
