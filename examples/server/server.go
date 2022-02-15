package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/voocel/ftf/network"
)

func main() {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&logrus.JSONFormatter{})
	s := network.NewServer("0.0.0.0:1235", network.WithHeartbeat(5*time.Second), network.WithLogger(l))

	s.OnConnect(func(c *network.Conn) {
		l.Println("connected by: %v\n", c.GetClientIP())
		go input(c)
	})

	s.OnMessage(func(c *network.Conn, msg *network.Message) {
		fmt.Println(msg)
	})

	s.OnClose(func(c *network.Conn, err error) {
		l.Printf("closed by[%v]: %v\n", c.GetClientIP(), err)
		ch := make(chan struct{}, 1)
		ch <- struct{}{}
		c.SetExtraMap(c.GetClientIP().String(), ch)
	})

	s.Start()
}

func input(c *network.Conn) {
	for {
		ch := c.GetExtraMap(c.GetClientIP().String())
		if ch != nil {
			cc := ch.(chan struct{})
			select {
			case <-cc:
				return
			default:
			}
		}

		reader := bufio.NewReader(os.Stdin)
		s, _ := reader.ReadString('\n')
		s = strings.Trim(s, "\n")
		c.SendBytes(network.Heartbeat, []byte(s))
	}
}
