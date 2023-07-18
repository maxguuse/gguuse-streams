package commands

import (
	"fmt"
	"log"

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

func (c *stopAnnouncementCommand) GetAnswer() (string, error) {
	if len(c.cmdArgs) != 1 {
		return "Usage: !stopannouncement <id>", nil
	}

	ann, isExists := repositories.Announcements.GetAnnouncement(c.cmdArgs[0])
	if !isExists {
		return fmt.Sprintf("Announcement with ID: %s doesn't exist", c.cmdArgs[0]), nil
	}

	repositories.Announcements.RemoveAnnouncement(ann.Id)
	err := repositories.Announcements.SaveAnnouncements()
	if err != nil {
		log.Fatalf("Error occured while saving announcements: %s", err)
	}

	return fmt.Sprintf("Announcement with ID: %s has been deleted", c.cmdArgs[0]), nil
}
