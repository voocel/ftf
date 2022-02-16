package main

import (
	"bufio"
	"fmt"
	"io"
	//"go/scanner"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/voocel/ftf/network"
)

type Client struct {
	conn net.Conn
	p    network.Protocol
	rch  chan *network.Message
	wch  chan *network.Message
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1235")
	if err != nil {
		panic(err)
	}

	c := &Client{
		conn: conn,
		p:    network.NewDefaultProtocol(),
		rch:  make(chan *network.Message, 1024),
		wch:  make(chan *network.Message, 1024),
	}

	go c.read(conn)
	go c.write(conn)
	go c.input()

	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-q
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func (c *Client) read(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		msg, err := c.p.Unpack(reader)
		if err != nil {
			fmt.Printf("read err: %v\n", err)
			return
		}
		fmt.Printf("receive: %v\n", msg)
	}
}

func (c *Client) write(conn net.Conn) {
	defer conn.Close()
	for {
		select {
		case msg := <-c.wch:
			b, err := c.p.Pack(msg)
			if err != nil {
				fmt.Printf("write pack err: %v\n")
				continue
			}
			c.conn.Write(b)
		}
	}
}

func (c *Client) input() {
	for {
		reader := bufio.NewReader(os.Stdin)
		s, _ := reader.ReadString('\n')
		s = strings.Trim(s, "\n")
		msg := network.NewMessage(network.Heartbeat, []byte(s))
		c.wch <- msg
	}
}

func tcpToHttp()  {
	con, err := net.Dial("tcp", "qq.com:80")
	if err != nil {
		panic(err)
	}
	defer con.Close()
	//fmt.Fprintf(con, "GET / HTTP/1.1\r\n\r\n")
	con.Write([]byte("GET / HTTP/1.1\\r\\n\\r\\n"))

	var sb strings.Builder
	buf := make([]byte, 256)
	for {
		n, err := con.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		sb.Write(buf[:n])
	}
	fmt.Printf("total response size:%d, response: %s\n", sb.Len(), sb.String())
}