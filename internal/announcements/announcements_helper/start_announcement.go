package announcements_helper

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/maxguuse/gguuse-streams/internal/dataaccess"
	"github.com/maxguuse/gguuse-streams/tools"
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

		broadcasterId := tools.GetUserId(channel)
		moderatorId := tools.GetUserId(os.Getenv("NICK"))

		sendAnnouncement(broadcasterId, moderatorId, ann.Text)

		log.Printf("Sent announcement '%s'", annId)

		time.Sleep(ann.RepetitionInterval)
	}
}


func sendAnnouncement(broadcasterId, moderatorId string, announcement string) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/chat/announcements?broadcaster_id=%s&moderator_id=%s", broadcasterId, moderatorId)

	formattedAnnouncement := fmt.Sprintf(`{"message":"%s","color":"primary"}`, announcement)
	jsonAnnouncement := []byte(formattedAnnouncement)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonAnnouncement))

	oauth := os.Getenv("AUTH_TOKEN")
	bearer := tools.BuildBearerFromOAuth(oauth)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Client-Id", os.Getenv("CLIENT_ID"))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending get users request")
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Error sending announcement")
	}
}
