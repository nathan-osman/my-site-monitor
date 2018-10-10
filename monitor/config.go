package monitor

import (
	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/nathan-osman/my-site-monitor/notifier"
)

// Config stores information for the monitor.
type Config struct {
	Conn     *db.Conn
	Notifier *notifier.Notifier
}
