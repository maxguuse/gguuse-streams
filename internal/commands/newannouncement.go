package commands

import (
	"fmt"
	"log"
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

func (c *newAnnouncementCommand) GetAnswer() (string, error) {
	if len(c.cmdArgs) < 3 {
		return "Usage: !newannouncement <id> <repetition_interval> <message>", nil
	}

	annId := c.cmdArgs[0]
	rawAnnRepTime, err := strconv.Atoi(c.cmdArgs[1])
	if err != nil {
		return "", err
	}
	annRepTime := time.Duration(rawAnnRepTime) * time.Minute
	annText := strings.Join(c.cmdArgs[2:], " ")

	ann := announcements.NewAnnouncement(annId, annRepTime, annText)

	repositories.Announcements.AddAnnouncement(*ann)
	err = repositories.Announcements.SaveAnnouncements()
	if err != nil {
		log.Fatalf("Error occured while saving announcements: %s", err)
	}

	go announcements_helper.StartAnnouncement(annId)

	return fmt.Sprintf("Created new announcement with ID: %s", annId), nil
}
