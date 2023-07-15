package announcments

import "time"

type Announcment struct {
	Id                 string        `json:"id"`
	RepetitionInterval time.Duration `json:"repetitionInterval"`
	Text               string        `json:"text"`
}

func NewAnnouncment(
	id string,
	repetitionInterval time.Duration,
	text string,
) *Announcment {
	return &Announcment{
		Id:                 id,
		RepetitionInterval: repetitionInterval,
		Text:               text,
	}
}

func (a *Announcment) GetID() string {
	return a.Id
}

func (a *Announcment) GetText() string {
	return a.Text
}

func (a *Announcment) GetRepTime() time.Duration {
	return a.RepetitionInterval
}
