package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	ip        string
	port      int32
	logger    *log.Logger
	tlsConf   *tls.Config
	onConnect func(c *Client)
	onMessage func(c *Client, msg string)
	onClose   func(c *Client, err error)
}

func NewServer(ip string, port int32) *Server {
	return &Server{
		ip:        ip,
		port:      port,
		onConnect: func(c *Client) {},
		onMessage: func(c *Client, msg string) {},
		onClose:   func(c *Client, err error) {},
		logger:    log.New(os.Stderr, "[server]", log.LstdFlags),
	}
}

func (s *Server) log(format string, v ...interface{}) {
	s.logger.Printf(format, v...)
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	s.log("server start success! %s:%d", s.ip, s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.log("accept connection err: %v", err)
			continue
		}
		c := Client{
			Conn:     conn,
			server:   s,
			clientIP: conn.RemoteAddr(),
		}
		go c.process()
	}
}

func (s *Server) OnConnect(callback func(c *Client)) {
	s.onConnect = callback
}

func (s *Server) OnMessage(callback func(c *Client, msg string)) {
	s.onMessage = callback
}

func (s *Server) OnClose(callback func(c *Client, err error)) {
	s.onClose = callback
}

type Client struct {
	Conn     net.Conn
	server   *Server
	clientIP net.Addr
}

func (c *Client) process() {
	c.server.onConnect(c)
	reader := bufio.NewReader(c.Conn)
	for {
		msg, err := Decode(reader)
		if err != nil {
			c.Conn.Close()
			c.server.onClose(c, err)
			return
		}
		c.server.onMessage(c, msg)
	}
}

func (c *Client) Send(msg string) error {
	return c.SendBytes([]byte(msg))
}

func (c *Client) SendBytes(b []byte) error {
	_, err := c.Conn.Write(b)
	if err != nil {
		c.Conn.Close()
		c.server.onClose(c, err)
	}
	return err
}

func (c *Client) GetClientIP() net.Addr {
	return c.clientIP
}
