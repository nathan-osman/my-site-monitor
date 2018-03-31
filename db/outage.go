package db

import (
	"time"
)

// Outage represents a period of time during which a site was inaccessible.
type Outage struct {
	ID int64

	// Start and ending time of the outage (may be null if ongoing)
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time

	// Information about the error
	Description string `gorm:"not null"`

	// Site that is experiencing the outage
	Site   *Site `gorm:"ForeignKey:SiteID"`
	SiteID int64 `sql:"type:int REFERENCES sites(id) ON DELETE CASCADE"`
}
