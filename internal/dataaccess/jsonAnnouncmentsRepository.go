package dataaccess

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/maxguuse/gguuse-streams/internal/announcments"
)

type jsonAnnouncmentsRepository struct {
	anns map[string]announcments.Announcment
	file string
}

func NewJsonAnnouncmentsRepository(channel string) *jsonAnnouncmentsRepository {
	return &jsonAnnouncmentsRepository{
		anns: make(map[string]announcments.Announcment),
		file: fmt.Sprintf("%s_announcments.json", channel),
	}
}

func (c *jsonAnnouncmentsRepository) GetAnnouncment(id string) (*announcments.Announcment, bool) {
	ann, isExists := c.anns[id]

	if !isExists {
		return nil, false
	}

	return &ann, true
}

func (c *jsonAnnouncmentsRepository) AddAnnouncment(ann announcments.Announcment) {
	c.anns[ann.GetID()] = ann
}

func (c *jsonAnnouncmentsRepository) RemoveAnnouncment(id string) {
	delete(c.anns, id)
}

func (c *jsonAnnouncmentsRepository) LoadAnnouncments() (err error) {
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	annsFile, err := os.OpenFile(c.file, os.O_RDONLY|os.O_CREATE, 0644)
	c.SaveAnnouncments()
	byteValue, err := io.ReadAll(annsFile)
	err = json.Unmarshal(byteValue, &c.anns)

	if err != nil {
		return err
	}
	defer annsFile.Close()
	return nil
}

func (c *jsonAnnouncmentsRepository) SaveAnnouncments() (err error) {
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
