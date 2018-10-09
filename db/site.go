package db

import (
	"encoding/json"
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

const (
	iso8601 = "2006-01-02T15:04:05-0700"
)

// Site represents a website being monitored.
type Site struct {
	ID int64 `json:"id"`

	// URL to poll and a descriptive name for the site
	URL  string `gorm:"not null" json:"url"`
	Name string `gorm:"not null" json:"name"`

	// Poll interval and the time of last and next poll
	PollInterval int64     `gorm:"not null" json:"poll-interval"`
	LastPoll     time.Time `gorm:"not null" json:"last-poll"`
	NextPoll     time.Time `gorm:"not null" json:"next-poll"`

	// Current status of the site and the time since it last changed
	Status     int       `gorm:"not null" json:"status"`
	StatusTime time.Time `json:"status-time"`

	// User that created the site
	User   *User `gorm:"ForeignKey:UserID" json:"-"`
	UserID int64 `sql:"type:int REFERENCES users(id) ON DELETE CASCADE" json:"user-id"`
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

// MarshalJSON ensures that dates are converted to ISO 8601.
func (s *Site) MarshalJSON() ([]byte, error) {
	type Alias Site
	return json.Marshal(&struct {
		*Alias
		LastPoll   string `json:"last-poll"`
		NextPoll   string `json:"next-poll"`
		StatusTime string `json:"status-time"`
	}{
		Alias:      (*Alias)(s),
		LastPoll:   s.LastPoll.Format(iso8601),
		NextPoll:   s.NextPoll.Format(iso8601),
		StatusTime: s.StatusTime.Format(iso8601),
	})
}
