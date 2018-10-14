package db

import (
	"strconv"
	"time"
)

// Outage represents a period of time during which a site was inaccessible.
type Outage struct {
	ID int64 `json:"id"`

	// Start and ending time of the outage (may be 0 if ongoing)
	StartTime time.Time  `gorm:"not null" json:"start-time"`
	EndTime   *time.Time `json:"end-time"`

	// Whether the two notifications have been sent yet or not
	StartNotificationSent bool `gorm:"not null" json:"-"`
	EndNotificationSent   bool `gorm:"not null" json:"-"`

	// Information about the error
	Description string `gorm:"not null" json:"description"`

	// Site that is experiencing the outage
	Site   *Site `gorm:"ForeignKey:SiteID" json:"-"`
	SiteID int64 `sql:"type:int REFERENCES sites(id) ON DELETE CASCADE" json:"site-id"`
}

// GetName retrieves a descriptive name for groups of outages
func (o *Outage) GetName() string {
	return "outages"
}

// GetID returns the unique identifier for the outage.
func (o *Outage) GetID() string {
	return strconv.FormatInt(o.ID, 10)
}

// SetID sets the unique identifier for the outage.
func (o *Outage) SetID(id string) error {
	o.ID, _ = strconv.ParseInt(id, 10, 64)
	return nil
}
