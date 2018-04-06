package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "msm"
	app.Usage = "My Site Monitor"
	app.Flags = flags
	app.Commands = []cli.Command{
		createadminCommand,
		runCommand,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err.Error())
	}
}
