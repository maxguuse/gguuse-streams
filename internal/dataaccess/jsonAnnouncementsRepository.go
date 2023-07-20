package dataaccess

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	twitch_config "github.com/maxguuse/gguuse-streams/configs/twitch"

	"github.com/maxguuse/gguuse-streams/internal/announcements"
	"golang.org/x/exp/maps"
)

type jsonAnnouncementsRepository struct {
	anns map[string]announcements.Announcement
	file string
}

func NewJsonAnnouncementsRepository() *jsonAnnouncementsRepository {
	if os.Getenv("GO_ENV") == "Production" {
		return &jsonAnnouncementsRepository{
			anns: make(map[string]announcements.Announcement),
			file: fmt.Sprintf("/home/gguuse/go/src/gguuse-streams/json_announcements/%s_announcements.json", twitch_config.Channel),
		}
	} else {
		return &jsonAnnouncementsRepository{
			anns: make(map[string]announcements.Announcement),
			file: fmt.Sprintf("json_announcements/%s_announcements.json", twitch_config.Channel),
		}
	}
}

func (c *jsonAnnouncementsRepository) GetAnnouncement(id string) (*announcements.Announcement, bool) {
	ann, isExists := c.anns[id]

	if !isExists {
		return nil, false
	}

	return &ann, true
}

func (c *jsonAnnouncementsRepository) GetIds() []string {
	return maps.Keys(c.anns)
}

func (c *jsonAnnouncementsRepository) AddAnnouncement(ann announcements.Announcement) {
	c.anns[ann.Id] = ann
}

func (c *jsonAnnouncementsRepository) RemoveAnnouncement(id string) {
	delete(c.anns, id)
}

func (c *jsonAnnouncementsRepository) LoadAnnouncements() (err error) {
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	annsFile, err := os.OpenFile(c.file, os.O_RDONLY|os.O_CREATE, 0644)
	byteValue, err := io.ReadAll(annsFile)

	if !json.Valid(byteValue) {
		c.SaveAnnouncements()
		byteValue, _ = io.ReadAll(annsFile)
	}

	err = json.Unmarshal(byteValue, &c.anns)

	if err != nil {
		return err
	}
	defer annsFile.Close()
	return nil
}

func (c *jsonAnnouncementsRepository) SaveAnnouncements() (err error) {
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	annsFile, err := os.Create(c.file)
	byteValue, err := json.Marshal(c.anns)
	_, err = annsFile.Write(byteValue)

	if err != nil {
		return err
	}
	defer annsFile.Close()
	return nil
}
