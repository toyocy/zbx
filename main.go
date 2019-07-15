package main

import (
	"log"
	"os"

	"github.com/toyocy/zbx/commands"
	"github.com/urfave/cli"
)

const version = "0.0.1"

type subCmd interface {
	Run(*cli.Context) error
}

func main() {
	app := cli.NewApp()
	app.Name = "zbx"
	app.Version = version
	app.Usage = "Command Line Interface for Zabbix Server"

	app.Commands = []cli.Command{
		commands.Login(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
