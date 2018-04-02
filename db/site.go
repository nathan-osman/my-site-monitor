package db

import (
	"strconv"
	"time"
)

const (
	// StatusUnknown indicates that the status of the site is not yet known.
	StatusUnknown = 0
	// StatusUp indicates that the website is responding to requests.
	StatusUp = 1
	// StatusDown indicates that the website is not responding to requests.
	StatusDown = 2
)

// Site represents a website being monitored.
type Site struct {
	ID int64

	// URL to poll and a descriptive name for the site
	URL  string `gorm:"not null"`
	Name string `gorm:"not null"`

	// Poll interval and the time since the last update
	PollInterval int64 `gorm:"not null"`
	PollTime     time.Time

	// Current status of the site and the time since it last changed
	Status     int `gorm:"not null"`
	StatusTime time.Time

	// User that created the site
	User   *User `gorm:"ForeignKey:UserID"`
	UserID int64 `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
}

// GetName retrieves a descriptive name for groups of sites.
func (s *Site) GetName() string {
	return "sites"
}

// GetID returns the unique identifier for the site.
func (s *Site) GetID() string {
	return strconv.FormatInt(s.ID, 10)
}

// SetID sets the unique identifier for the site.
func (s *Site) SetID(id string) error {
	s.ID, _ = strconv.ParseInt(id, 10, 64)
	return nil
}
