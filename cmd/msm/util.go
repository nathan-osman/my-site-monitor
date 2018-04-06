package main

import (
	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/urfave/cli"
)

func initDB(c *cli.Context) (conn *db.Conn, err error) {
	conn, err = db.New(&db.Config{
		Host:     c.GlobalString("db-host"),
		Port:     c.GlobalInt("db-port"),
		Name:     c.GlobalString("db-name"),
		User:     c.GlobalString("db-user"),
		Password: c.GlobalString("db-password"),
	})
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			conn.Close()
			conn = nil
		}
	}()
	if err = conn.Migrate(); err != nil {
		return
	}
	return
}
