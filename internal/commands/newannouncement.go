package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/maxguuse/gguuse-streams/configs/repositories"
	"github.com/maxguuse/gguuse-streams/internal/announcements"
	"github.com/maxguuse/gguuse-streams/internal/announcements/announcements_helper"
)

type newAnnouncementCommand struct {
	cmdArgs []string
}

func NewNewAnnouncementCommand(cmdArgs []string) *newAnnouncementCommand {
	return &newAnnouncementCommand{
		cmdArgs: cmdArgs,
	}
}

func (c *newAnnouncementCommand) GetAnswer() string {
	if len(c.cmdArgs) < 3 {
		return "Usage: !newannouncement <id> <repetition_interval> <message>"
	}

	annId := c.cmdArgs[0]
	rawAnnRepTime, err := strconv.Atoi(c.cmdArgs[1])
	if err != nil {
		return ""
	}
	annRepTime := time.Duration(rawAnnRepTime) * time.Minute
	annText := strings.Join(c.cmdArgs[2:], " ")

	ann := announcements.NewAnnouncement(annId, annRepTime, annText)

	repositories.Announcements.AddAnnouncement(*ann)
	repositories.Announcements.SaveAnnouncements()

	go announcements_helper.StartAnnouncement(annId)

	return ""
}
