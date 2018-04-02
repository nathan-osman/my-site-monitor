package server

import (
	"github.com/nathan-osman/my-site-monitor/db"
)

// Config stores the configuration for the embedded web server.
type Config struct {
	Addr string
	Conn *db.Conn
}
