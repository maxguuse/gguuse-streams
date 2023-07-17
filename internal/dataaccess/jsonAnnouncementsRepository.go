package dataaccess

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/maxguuse/gguuse-streams/internal/announcements"
)

type jsonAnnouncementsRepository struct {
	anns map[string]announcements.Announcement
	file string
}

func NewJsonAnnouncementsRepository(channel string) *jsonAnnouncementsRepository {
	return &jsonAnnouncementsRepository{
		anns: make(map[string]announcements.Announcement),
		file: fmt.Sprintf("json_announcements/%s_announcements.json", channel),
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
	keys := make([]string, 0, len(c.anns))
	for k := range c.anns {
		keys = append(keys, k)
	}
	return keys
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
