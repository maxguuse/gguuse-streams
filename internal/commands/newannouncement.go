package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/maxguuse/gguuse-streams/internal/announcements"
	"github.com/maxguuse/gguuse-streams/internal/announcements/announcements_helper"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
)

type newAnnouncementCommand struct {
	anns    dataaccess.AnnouncementsRepository
	cmdArgs []string
	channel string
}

func NewNewAnnouncementCommand(
	anns dataaccess.AnnouncementsRepository,
	cmdArgs []string,
	channel string,
) *newAnnouncementCommand {
	return &newAnnouncementCommand{
		anns:    anns,
		cmdArgs: cmdArgs,
		channel: channel,
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

	c.anns.AddAnnouncement(*ann)
	c.anns.SaveAnnouncements()

	go announcements_helper.StartAnnouncement(annId, c.anns, c.channel)

	return ""
}
