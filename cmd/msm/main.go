package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/howeyc/gopass"
	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/nathan-osman/my-site-monitor/monitor"
	"github.com/nathan-osman/my-site-monitor/server"
	"github.com/sirupsen/logrus"
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
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
			Usage:  "enable debug logging",
		},
		cli.StringFlag{
			Name:   "secret-key",
			EnvVar: "SECRET_KEY",
			Usage:  "secret key for sessions",
		},
		cli.StringFlag{
			Name:   "server-addr",
			Value:  ":8000",
			EnvVar: "SERVER_ADDR",
			Usage:  "server address",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "createuser",
			Usage: "create a new user",
			Action: func(c *cli.Context) error {

				// Connect to the database
				conn, err := db.New(&db.Config{
					Host:     c.GlobalString("db-host"),
					Port:     c.GlobalInt("db-port"),
					Name:     c.GlobalString("db-name"),
					User:     c.GlobalString("db-user"),
					Password: c.GlobalString("db-password"),
				})
				if err != nil {
					return err
				}
				defer conn.Close()

				// Prompt for the username
				var username string
				fmt.Print("Username? ")
				fmt.Scanln(&username)

				// Prompt for the password, hiding the input
				fmt.Print("Password? ")
				p, err := gopass.GetPasswd()
				if err != nil {
					return err
				}

				// Create a user and set their password
				u := &db.User{
					Username: username,
				}
				if err := u.SetPassword(string(p)); err != nil {
					return err
				}

				// Save the user to the database
				if err := conn.Save(u).Error; err != nil {
					return err
				}

				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {

		// Enable debug logging if requested
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

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

		// Perform all database migrations
		if err = conn.Migrate(); err != nil {
			return err
		}

		// Create the monitor
		m := monitor.New(&monitor.Config{
			Conn: conn,
		})
		defer m.Close()

		// Start the server
		s, err := server.New(&server.Config{
			Addr:      c.String("server-addr"),
			Conn:      conn,
			Monitor:   m,
			SecretKey: c.String("secret-key"),
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
