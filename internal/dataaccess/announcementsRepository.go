package dataaccess

import "github.com/maxguuse/gguuse-streams/internal/announcements"

type AnnouncementsRepository interface {
	GetAnnouncement(id string) (*announcements.Announcement, bool)
	GetIds() []string
	AddAnnouncement(announcement announcements.Announcement)
	RemoveAnnouncement(id string)

	LoadAnnouncements() error
	SaveAnnouncements() error
}
