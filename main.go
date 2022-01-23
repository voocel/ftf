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
		Commands: []*cli.Command{
			{
				Name:   "send",
				Usage:  "send server",
				Action: send,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "ip", Usage: "输入目标地址", Value: defaultAddr},
				},
			},
			{
				Name:   "receive",
				Usage:  "receive server",
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
