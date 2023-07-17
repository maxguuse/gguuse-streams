package commands

import (
	"fmt"

	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
)

type stopAnnouncementCommand struct {
	anns    dataaccess.AnnouncementsRepository
	cmdArgs []string
}

func NewStopAnnouncementCommand(
	anns dataaccess.AnnouncementsRepository,
	cmdArgs []string,
) *stopAnnouncementCommand {
	return &stopAnnouncementCommand{
		anns:    anns,
		cmdArgs: cmdArgs,
	}
}

func (c *stopAnnouncementCommand) GetAnswer() string {
	if len(c.cmdArgs) != 1 {
		return "Usage: !stopannouncement <id>"
	}

	ann, isExists := c.anns.GetAnnouncement(c.cmdArgs[0])
	if !isExists {
		return fmt.Sprintf("Announcement with ID: %s doesn't exist", c.cmdArgs[0])
	}

	c.anns.RemoveAnnouncement(ann.Id)
	c.anns.SaveAnnouncements()

	return fmt.Sprintf("Announcement with ID: %s has been deleted", c.cmdArgs[0])
}
