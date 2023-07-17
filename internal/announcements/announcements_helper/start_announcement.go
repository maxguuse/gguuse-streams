package announcements_helper

import (
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
)

func StartAnnouncement(
	annId string,
	anns dataaccess.AnnouncementsRepository,
	client *twitch.Client,
	channel string,
) {
	log.Printf("Started announcement '%s'", annId)
	for {
		ann, isExists := anns.GetAnnouncement(annId)
		if !isExists {
			break
		}

		// TODO: Replace this with HTTP request so announcement show like announcement
		client.Say(channel, ann.Text)

		log.Printf("Sent announcement '%s'", annId)

		time.Sleep(ann.RepetitionInterval)
	}
}
