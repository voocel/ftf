package main

import (
	"f2f/server"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

var FileName string

func receive(context *cli.Context) (err error) {
	s := server.NewServer("0.0.0.0", 1234)
	s.OnConnect(func(c *server.Client) {
		fmt.Printf("connect by: %v\n", c.GetClientIP())
	})

	s.OnMessage(func(c *server.Client, msg string) {
		if strings.Index(msg, "ack") == 0 {
			err := c.SendBytes([]byte("ok"))
			if err != nil {
				return
			}
			FileName = "r-" + msg[3:]
		}
		saveFile(msg)
	})

	s.OnClose(func(c *server.Client, err error) {
		fmt.Printf("closed by: %v\n", c.GetClientIP())
	})

	s.Start()
	return
}

func saveFile(msg string) {
	// 路径穿越检验
	if strings.Contains(FileName, "..") {
		return
	}
	file, err := os.Create(FileName)
	if err != nil {
		log.Fatalln("create file err: ", err)
		return
	}
	defer file.Close()

	file.Write([]byte(msg))
}

func checkIllegal(cmdName string) bool {
	if strings.Contains(cmdName, "&") || strings.Contains(cmdName, "|") ||
		strings.Contains(cmdName, ";") || strings.Contains(cmdName, "$") ||
		strings.Contains(cmdName, "'") || strings.Contains(cmdName, "`") ||
		strings.Contains(cmdName, "(") || strings.Contains(cmdName, ")") ||
		strings.Contains(cmdName, "\"") {
		return true
	}
	return false
}
