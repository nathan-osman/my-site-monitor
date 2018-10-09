package notifier

import (
	"github.com/nathan-osman/my-site-monitor/db"
)

// Config stores information for the notifier.
type Config struct {
	Conn           *db.Conn
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}
