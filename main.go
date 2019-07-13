package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const version = "0.0.1"

type subCmd interface {
	Run(*cli.Context) error
}

var cmdList = []cli.Command{}

func main() {
	app := cli.NewApp()
	app.Name = "zbx"
	app.Version = version
	app.Usage = "Command Line Interface for Zabbix Server"

	app.Commands = cmdList

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context, sub subCmd) error {
	return sub.Run(c)
}
