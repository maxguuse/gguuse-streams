package announcements_helper

import (
	"log"
	"os"
	"time"

	twitch_config "github.com/maxguuse/gguuse-streams/configs/twitch"
	"github.com/nicklaw5/helix/v2"

	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
)

func StartAnnouncement(
	annId string,
	anns dataaccess.AnnouncementsRepository,
	channel string,
) {
	usersResp, err := twitch_config.ApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{channel, os.Getenv("NICK")},
	})
	if err != nil {
		log.Printf("Error fetching broadcaster id: %s", err)
		return
	}

	var broadcasterId, moderatorId string

	for _, v := range usersResp.Data.Users {
		if v.Login == channel {
			broadcasterId = v.ID
		} else if v.Login == os.Getenv("NICK") {
			moderatorId = v.ID
		}
	}

	if channel == os.Getenv("NICK") {
		moderatorId = broadcasterId
	}

	sendChatAnnouncementParams := &helix.SendChatAnnouncementParams{
		BroadcasterID: broadcasterId,
		ModeratorID:   moderatorId,
		Color:         "primary",
	}

	log.Printf("Started announcement '%s'", annId)
	for {
		ann, isExists := anns.GetAnnouncement(annId)
		if !isExists {
			log.Printf("Announcement '%s' doesn't exist", annId)
			return
		}

		sendChatAnnouncementParams.Message = ann.Text
		_, err = twitch_config.ApiClient.SendChatAnnouncement(sendChatAnnouncementParams)

		if err != nil {
			log.Printf("Error sending announcement '%s', error: %s", ann.Id, err)
			return
		}

		log.Printf("Sent announcement '%s'", annId)

		time.Sleep(ann.RepetitionInterval)
	}
}
