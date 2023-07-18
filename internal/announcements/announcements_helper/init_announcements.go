package announcements_helper

import (
	"log"

	"github.com/maxguuse/gguuse-streams/configs/repositories"
)

func InitAnnouncements() {
	log.Println("Started initialization of announcements")
	ids := repositories.Announcements.GetIds()
	for _, annId := range ids {
		go StartAnnouncement(annId)
	}
	log.Println("All announcements initialized")
}
