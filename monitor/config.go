package monitor

import (
	"github.com/nathan-osman/my-site-monitor/db"
)

// Config stores information for the monitor.
type Config struct {
	Conn *db.Conn
}
