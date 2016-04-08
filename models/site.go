package models

import (
	"github.com/jinzhu/gorm"
)

// Site represents an individual website being monitored.
type Site struct {
	gorm.Model
	Url     string
	Timeout int
}
