package models

const (
	LevelInfo = iota
	LevelWarning
	LevelError
)

// Alert represents an event notification. Usually this indicates that the
// status of a site has changed.
type Alert struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	Site      Site
	SiteID    uint
	Level     uint
}
