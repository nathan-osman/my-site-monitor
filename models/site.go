package models

const (
	StatusUnknown = iota
	StatusOK
	StatusWarning
	StatusError
)

// Site represents an individual website being monitored.
type Site struct {
	ID uint `gorm:"primary_key"`

	URL    string // Full URI of site
	Status uint   // Current status of the site

	IntervalMS uint // Time between pings
	WarningMS  uint // Minimum response time to raise a warning
	ErrorMS    uint // Minimum response time to raise an error
}
