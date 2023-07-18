package announcements_helper

import (
	"log"

	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
)

func InitAnnouncements(
	anns dataaccess.AnnouncementsRepository,
) {
	log.Println("Started initialization of announcements")
	ids := anns.GetIds()
	for i := 0; i < len(ids); i++ {
		go StartAnnouncement(ids[i], anns)
	}
	log.Println("All announcements initialized")
}
