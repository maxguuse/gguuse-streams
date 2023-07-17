package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)
func BuildBearerFromOAuth(oauth string) (bearer string) {
	return fmt.Sprintf("Bearer %s", oauth[6:])
}


func GetUserId(login string) (userId string) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", login)
	req, _ := http.NewRequest("GET", url, nil)

	oauth := os.Getenv("AUTH_TOKEN")
	bearer := BuildBearerFromOAuth(oauth)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Client-Id", os.Getenv("CLIENT_ID"))

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending get users request")
		return ""
	}

	respByteValue, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if !json.Valid(respByteValue) {
		log.Println("Invalid json from get users request")
		return ""
	}
	respMap := make(map[string]any)
	json.Unmarshal(respByteValue, &respMap)

	return respMap["data"].([]any)[0].(map[string]any)["id"].(string)
}
