package announcements

import "time"

type Announcement struct {
	Id                 string        `json:"id"`
	RepetitionInterval time.Duration `json:"repetitionInterval"`
	Text               string        `json:"text"`
}

func NewAnnouncement(
	id string,
	repetitionInterval time.Duration,
	text string,
) *Announcement {
	return &Announcement{
		Id:                 id,
		RepetitionInterval: repetitionInterval,
		Text:               text,
	}
}
