package commands

import (
	"fmt"

	"github.com/maxguuse/gguuse-streams/configs/repositories"
)

type stopAnnouncementCommand struct {
	cmdArgs []string
}

func NewStopAnnouncementCommand(cmdArgs []string) *stopAnnouncementCommand {
	return &stopAnnouncementCommand{
		cmdArgs: cmdArgs,
	}
}

func (c *stopAnnouncementCommand) GetAnswer() string {
	if len(c.cmdArgs) != 1 {
		return "Usage: !stopannouncement <id>"
	}

	ann, isExists := repositories.Announcements.GetAnnouncement(c.cmdArgs[0])
	if !isExists {
		return fmt.Sprintf("Announcement with ID: %s doesn't exist", c.cmdArgs[0])
	}

	repositories.Announcements.RemoveAnnouncement(ann.Id)
	repositories.Announcements.SaveAnnouncements()

	return fmt.Sprintf("Announcement with ID: %s has been deleted", c.cmdArgs[0])
}
