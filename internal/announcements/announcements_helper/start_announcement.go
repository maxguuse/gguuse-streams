package announcements_helper

import (
	"log"
	"os"
	"time"

	"github.com/maxguuse/gguuse-streams/configs/repositories"
	twitch_config "github.com/maxguuse/gguuse-streams/configs/twitch"

	"github.com/nicklaw5/helix/v2"
)

func StartAnnouncement(annId string) {
	usersResp, err := twitch_config.ApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{twitch_config.Channel, os.Getenv("NICK")},
	})
	if err != nil {
		log.Printf("Error fetching broadcaster id: %s", err)
		return
	}

	var broadcasterId, moderatorId string

	for _, v := range usersResp.Data.Users {
		if v.Login == twitch_config.Channel {
			broadcasterId = v.ID
		} else if v.Login == os.Getenv("NICK") {
			moderatorId = v.ID
		}
	}

	if twitch_config.Channel == os.Getenv("NICK") {
		moderatorId = broadcasterId
	}

	sendChatAnnouncementParams := &helix.SendChatAnnouncementParams{
		BroadcasterID: broadcasterId,
		ModeratorID:   moderatorId,
		Color:         "primary",
	}

	log.Printf("Started announcement '%s'", annId)
	for {
		ann, isExists := repositories.Announcements.GetAnnouncement(annId)
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
