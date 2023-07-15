package dataaccess

import "github.com/maxguuse/gguuse-streams/internal/announcments"

type AnnouncmentsRepository interface {
	GetAnnouncment(id string) (*announcments.Announcment, bool)
	AddAnnouncment(announcment announcments.Announcment)
	RemoveAnnouncment(id string)

	SaveAnnouncments() error
	LoadAnnouncments() error
}
