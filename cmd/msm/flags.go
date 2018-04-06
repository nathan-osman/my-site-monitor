package main

import (
	"github.com/urfave/cli"
)

var flags = []cli.Flag{
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
