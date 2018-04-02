package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/nathan-osman/my-site-monitor/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "msm"
	app.Usage = "My Site Monitor"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db-filename",
			Value:  "msm.sqlite3",
			EnvVar: "DB_FILENAME",
			Usage:  "filename for database file",
		},
		cli.StringFlag{
			Name:   "server-addr",
			Value:  ":8000",
			EnvVar: "SERVER_ADDR",
			Usage:  "server address",
		},
	}
	app.Action = func(c *cli.Context) error {

		// Connect to the database
		conn, err := db.New(&db.Config{
			Filename: c.String("dir"),
		})
		if err != nil {
			return err
		}
		defer conn.Close()

		// Perform all pending migrations
		if err := conn.Migrate(); err != nil {
			return err
		}

		// Start the server
		s, err := server.New(&server.Config{
			Addr: c.String("addr"),
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
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err.Error())
	}
}
