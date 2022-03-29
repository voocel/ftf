package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/voocel/ftf/network"
)

var flog network.Logger

type Receive struct {
	addr     string
	logger   *log.Logger
	protocol network.Protocol
}

func init() {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	flog = l
}

func receiver(ctx *cli.Context) (err error) {
	r := NewReceive("0.0.0.0:1234")
	r.receive()
	return
}

func NewReceive(addr string) *Receive {
	return &Receive{
		addr:     addr,
		protocol: network.NewDefaultProtocol(),
	}
}

func (r *Receive) receive() {
	s := network.NewServer(r.addr, network.WithLogger(flog))
	s.OnConnect(func(c *network.Conn) {
		flog.Infof("connected by: %v\n", c.GetClientIP())
	})

	s.OnMessage(func(c *network.Conn, msg *network.Message) {
		if msg.GetCmd() == network.Ack {
			c.SendBytes(network.Ack, []byte("ok"))
			c.SetExtraMap("filename", string(msg.GetData()))
			return
		}
		saveFile(c, msg)
	})

	s.OnClose(func(c *network.Conn, err error) {
		flog.Infof("closed by[%v]: %v\n", c.GetClientIP(), err)
	})

	s.Start()
}

func saveFile(c *network.Conn, msg *network.Message) {
	fileName := c.GetExtraMap("filename").(string)
	if fileExists(fileName) {
		suffix := filepath.Ext(fileName)
		prefix := strings.TrimSuffix(fileName, suffix)
		fileName = fmt.Sprintf("%s_ftf%s", prefix, suffix)
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

	//w := bufio.NewWriter(file)
	//io.Copy(w, bytes.NewReader(msg.GetData()))
	//w.Flush()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			// todo canceled the download
			return
		default:
			_, err := io.CopyN(file, bytes.NewReader(msg.GetData()), 4096)
			if err != nil {
				if err == io.EOF {
					return
				} else {
					flog.Fatal(err)
				}
			}
		}
	}
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
