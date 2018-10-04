package server

import (
	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/nathan-osman/my-site-monitor/monitor"
)

// Config stores the configuration for the embedded web server.
type Config struct {
	Addr      string
	Conn      *db.Conn
	Monitor   *monitor.Monitor
	SecretKey string
}
