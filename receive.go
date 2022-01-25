package main

import (
	"fmt"
	"ftf/network"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func receive(ctx *cli.Context) (err error) {
	s := network.NewServer("0.0.0.0", 1234)
	s.OnConnect(func(c *network.Client) {
		fmt.Printf("connect by: %v\n", c.GetClientIP())
	})

	s.OnMessage(func(c *network.Client, msg *network.Message) {
		if msg.GetCmd() == network.Ack {
			msg := network.NewMessage(network.Ack, []byte("ok"))
			c.SendMessage(msg)
			c.ExtraMap["filename"] = string(msg.GetData())
			return
		}
		saveFile(c, msg)
	})

	s.OnClose(func(c *network.Client, err error) {
		fmt.Printf("closed by: %v\n", c.GetClientIP())
	})

	s.Start()
	return
}

func saveFile(c *network.Client, msg *network.Message) {
	fileName, ok := c.ExtraMap["filename"]
	if !ok {
		fmt.Printf("filename  not exists: %v\n", c.GetClientIP())
		return
	}
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
		log.Fatalln("create file err: ", err)
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
