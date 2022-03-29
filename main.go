package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	defaultAddr = "127.0.0.1:1234"
)

func main() {
	app := cli.NewApp()
	app.Usage = "The high-performance file transfer for Golang"
	app.Version = Version
	app.EnableBashCompletion = true
	app.CustomAppHelpTemplate = Summary()

	app.Action = baseCmd
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name:   "send",
			Usage:  "send network",
			Action: sender,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "ip", Usage: "Enter target address", Value: defaultAddr},
			},
		},
		{
			Name:   "receive",
			Usage:  "receive network",
			Action: receiver,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}
}

func baseCmd(c *cli.Context) (err error) {
	//if c.NArg() == 0 {
	//	print(Summary())
	//	return
	//}
	if err := runApp(c); err != nil {
		return cli.Exit(err, 1)
	}
	return
}

func runApp(c *cli.Context) (err error) {
	opts := &startOpts{}
	opts.addr = c.String("addr")
	opts.path = c.String("path")

	a := newApp(opts)
	service := c.Command.Name
	if len(service) != 0 {
		switch service {
		case Sender:
			a.isSender = true
			a.hasService = true
		case Receiver:
			a.isReceiver = true
			a.hasService = true
		}
	}
	err = a.Start()
	if err != nil {
		return err
	}

	return err
}

func readMessageFromStdin() ([]byte, error) {
	var message []byte
	s, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if s.Mode()&os.ModeNamedPipe != 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		message = bytes
	}

	return message, nil
}
