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
	app := &cli.App{
		CustomAppHelpTemplate: Summary(),
		Commands: []*cli.Command{
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
		},
		Flags: []cli.Flag{},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}
}
