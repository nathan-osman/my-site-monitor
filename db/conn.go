package db

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	// Use an SQLite database internally
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Conn represents a connection to the database.
type Conn struct {
	*gorm.DB
	log *logrus.Entry
}

// New creates a new database connection and initializes it.
func New(cfg *Config) (*Conn, error) {
	g, err := gorm.Open("sqlite3", cfg.Filename)
	if err != nil {
		return nil, err
	}
	c := &Conn{
		DB:  g,
		log: logrus.WithField("context", "db"),
	}
	return c, nil
}

// Migrate performs all database migrations.
func (c *Conn) Migrate() error {
	c.log.Info("performing migrations...")
	return c.AutoMigrate(
		&User{},
		&Site{},
		&Outage{},
	).Error
}
