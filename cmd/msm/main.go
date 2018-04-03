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
			Name:   "db-host",
			Value:  "postgres",
			EnvVar: "DB_HOST",
			Usage:  "PostgreSQL database host",
		},
		cli.IntFlag{
			Name:   "db-port",
			Value:  5432,
			EnvVar: "DB_PORT",
			Usage:  "PostgreSQL database port",
		},
		cli.StringFlag{
			Name:   "db-name",
			Value:  "postgres",
			EnvVar: "DB_NAME",
			Usage:  "PostgreSQL database name",
		},
		cli.StringFlag{
			Name:   "db-user",
			Value:  "postgres",
			EnvVar: "DB_NAME",
			Usage:  "PostgreSQL database user",
		},
		cli.StringFlag{
			Name:   "db-password",
			Value:  "postgres",
			EnvVar: "DB_PASSWORD",
			Usage:  "PostgreSQL database password",
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
			Host:     c.String("db-host"),
			Port:     c.Int("db-port"),
			Name:     c.String("db-name"),
			User:     c.String("db-user"),
			Password: c.String("db-password"),
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
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err.Error())
	}
}
