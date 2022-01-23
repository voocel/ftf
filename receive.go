package main

import (
	"fmt"
	"ftf/server"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func receive(context *cli.Context) (err error) {
	s := server.NewServer("0.0.0.0", 1234)
	s.OnConnect(func(c *server.Client) {
		fmt.Printf("connect by: %v\n", c.GetClientIP())
	})

	s.OnMessage(func(c *server.Client, msg string) {
		if strings.Index(msg, "ack") == 0 {
			err := c.Send("ok")
			if err != nil {
				return
			}
			c.ExtraMap["filename"] = msg[3:]
			return
		}
		saveFile(c, msg)
	})

	s.OnClose(func(c *server.Client, err error) {
		fmt.Printf("closed by: %v\n", c.GetClientIP())
	})

	s.Start()
	return
}

func saveFile(c *server.Client, msg string) {
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

	file.Write([]byte(msg))
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
