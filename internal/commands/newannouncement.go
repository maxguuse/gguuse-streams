package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/maxguuse/gguuse-streams/internal/announcements"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
)

type newAnnouncementCommand struct {
	anns    dataaccess.AnnouncementsRepository
	cmdArgs []string
	client  *twitch.Client
	channel string
}

func NewNewAnnouncementCommand(
	anns dataaccess.AnnouncementsRepository,
	cmdArgs []string,
	client *twitch.Client,
	channel string,
) *newAnnouncementCommand {
	return &newAnnouncementCommand{
		anns:    anns,
		cmdArgs: cmdArgs,
		client:  client,
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
	go c.startAnnouncement(annId)

	return ""
}

func (c *newAnnouncementCommand) startAnnouncement(annId string) {
	for {
		ann, isExists := c.anns.GetAnnouncement(annId)
		if !isExists {
			break
		}

		// TODO: Replace this with HTTP request so announcement show like announcement
		c.client.Say(c.channel, ann.Text)

		time.Sleep(ann.RepetitionInterval)
	}
}
