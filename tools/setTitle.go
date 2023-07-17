package tools

import (
	"fmt"
	"bytes"
	"net/http"
	"os"
	"log"
)

func SetTitle(broadcasterId, title string ){
	url := "https://api.twitch.tv/helix/channels?broadcaster_id=" + broadcasterId

	formattedTitle := fmt.Sprintf(`{"title":"%s"}`, title)
	jsonTitle := []byte(formattedTitle)

	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonTitle))

	oauth := os.Getenv("AUTH_TOKEN")
	bearer := BuildBearerFromOAuth(oauth)
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println("Error changing title")
	}
}