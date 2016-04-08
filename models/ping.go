package models

import (
	"time"
)

// Ping represents an attempt to reach the specified site.
type Ping struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	Site       Site
	SiteID     uint
	ResponseMS uint
}
