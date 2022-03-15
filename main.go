package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const (
	defaultAddr = "127.0.0.1:1234"
)

func main() {
	app := cli.NewApp()
	app.Usage = "The high-performance file transfer for Golang"
	app.Version = Version
	app.EnableBashCompletion = true
	//app.CustomAppHelpTemplate = Summary()

	app.Action = baseCmd
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name:   "send",
			Usage:  "send network",
			Action: send,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "ip", Usage: "Enter target address", Value: defaultAddr},
			},
		},
		{
			Name:   "receive",
			Usage:  "receive network",
			Action: receive,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}
}

func baseCmd(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		print(Summary())
		return
	}
	return
}
