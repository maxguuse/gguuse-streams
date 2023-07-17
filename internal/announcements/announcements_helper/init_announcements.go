package announcements_helper

import (
	"log"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
)

func InitAnnouncements(
	anns dataaccess.AnnouncementsRepository,
	client *twitch.Client,
	channel string,
) {
	log.Println("Started initialization of announcements")
	ids := anns.GetIds()
	for i := 0; i < len(ids); i++ {
		go StartAnnouncement(
			ids[i],
			anns,
			client,
			channel,
		)
	}
	log.Println("All announcements initialized")
}
