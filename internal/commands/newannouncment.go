package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/maxguuse/gguuse-streams/internal/announcments"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
)

type newAnnouncmentCommand struct {
	anns    dataaccess.AnnouncmentsRepository
	cmdArgs []string
	client  *twitch.Client
	channel string
}

func NewNewAnnouncmentCommand(
	anns dataaccess.AnnouncmentsRepository,
	cmdArgs []string,
	client *twitch.Client,
	channel string,
) *newAnnouncmentCommand {
	return &newAnnouncmentCommand{
		anns:    anns,
		cmdArgs: cmdArgs,
		client:  client,
		channel: channel,
	}
}

func (c *newAnnouncmentCommand) GetAnswer() string {
	if len(c.cmdArgs) < 3 {
		return "Usage: !newannouncment <id> <repetition_interval> <message>"
	}

	annId := c.cmdArgs[0]
	rawAnnRepTime, err := strconv.Atoi(c.cmdArgs[1])
	if err != nil {
		return ""
	}
	annRepTime := time.Duration(rawAnnRepTime) * time.Minute
	annText := strings.Join(c.cmdArgs[2:], " ")

	ann := announcments.NewAnnouncment(annId, annRepTime, annText)

	c.anns.AddAnnouncment(*ann)
	c.anns.SaveAnnouncments()
	go c.startAnnouncment(annId)

	return ""
}

func (c *newAnnouncmentCommand) startAnnouncment(annId string) {
	for {
		ann, isExists := c.anns.GetAnnouncment(annId)
		if !isExists {
			break
		}

		// TODO: Replace this with HTTP request so announcement show like announcement
		c.client.Say(c.channel, ann.GetText())

		time.Sleep(ann.GetRepTime())
	}
}
