package main

import (
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/urfave/cli"
)

var createadminCommand = cli.Command{
	Name:  "createadmin",
	Usage: "create an admin user",
	Action: func(c *cli.Context) error {

		// Connect to the database
		conn, err := initDB(c)
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
			IsAdmin:  true,
		}
		if err := u.SetPassword(string(p)); err != nil {
			return err
		}

		return nil
	},
}
