package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/voocel/ftf/network"
)

var flog network.Logger

func init() {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	flog = l
}

func receive(ctx *cli.Context) (err error) {
	s := network.NewServer("0.0.0.0:1234", network.WithLogger(flog))
	s.OnConnect(func(c *network.Client) {
		flog.Infof("connected by: %v\n", c.GetClientIP())
	})

	s.OnMessage(func(c *network.Client, msg *network.Message) {
		if msg.GetCmd() == network.Ack {
			c.SendBytes(network.Ack, []byte("ok"))
			c.SetExtraMap("filename", string(msg.GetData()))
			return
		}
		saveFile(c, msg)
	})

	s.OnClose(func(c *network.Client, err error) {
		flog.Infof("closed by[%v]: %v\n", c.GetClientIP(), err)
	})

	s.Start()
	return
}

func saveFile(c *network.Client, msg *network.Message) {
	fileName := c.GetExtraMap("filename").(string)
	if fileExists(fileName) {
		suffix := filepath.Ext(fileName)
		prefix := strings.TrimSuffix(fileName, suffix)
		fileName = fmt.Sprintf("%s(FTF)%s", prefix, suffix)
	}
	if strings.Contains(fileName, "..") {
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		flog.Errorf("create file[%s] err: %v", fileName, err)
		return
	}
	defer file.Close()

	file.Write(msg.GetData())
}

func checkIllegal(cmdName string) (err bool) {
	if strings.Contains(cmdName, "&") || strings.Contains(cmdName, "|") ||
		strings.Contains(cmdName, ";") || strings.Contains(cmdName, "$") ||
		strings.Contains(cmdName, "'") || strings.Contains(cmdName, "`") ||
		strings.Contains(cmdName, "(") || strings.Contains(cmdName, ")") ||
		strings.Contains(cmdName, "\"") {
		return true
	}
	return
}
