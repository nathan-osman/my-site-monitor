package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/my-site-monitor/server"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "run the application",
	Action: func(c *cli.Context) error {

		// Connect to the database
		conn, err := initDB(c)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Start the server
		s, err := server.New(&server.Config{
			Addr: c.GlobalString("server-addr"),
			Conn: conn,
		})
		if err != nil {
			return err
		}
		defer s.Close()

		// Wait for SIGINT or SIGTERM
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		return nil
	},
}
