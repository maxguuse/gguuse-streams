package announcements_helper

import (
	"log"
	"os"
	"strings"
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
		log.Fatalf("Couldn't make GetUsers request to Twitch API, error: %s", err)
	} else if usersResp.ErrorMessage != "" {
		log.Println(usersResp.ErrorMessage)
		return
	}

	var broadcasterId, moderatorId string

	for _, v := range usersResp.Data.Users {
		if strings.EqualFold(v.Login, twitch_config.Channel) {
			broadcasterId = v.ID
		} else if strings.EqualFold(v.Login, os.Getenv("NICK")) {
			moderatorId = v.ID
		}
	}

	if strings.EqualFold(twitch_config.Channel, os.Getenv("NICK")) {
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
		sendChatAnnouncementResp, err := twitch_config.ApiClient.SendChatAnnouncement(sendChatAnnouncementParams)

		if err != nil {
			log.Fatalf("Couldn't make SendChatAnnouncement request to Twitch API, error: %s", err)
		} else if sendChatAnnouncementResp.ErrorMessage != "" {
			log.Println(sendChatAnnouncementResp.ErrorMessage)
			return
		}

		log.Printf("Sent announcement '%s'", annId)

		time.Sleep(ann.RepetitionInterval)
	}
}
