package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
)

// Server defines parameters for running an TCP server
type Server struct {
	ip        string
	port      int32
	logger    *log.Logger
	tlsConf   *tls.Config
	onConnect func(c *Client)
	onMessage func(c *Client, msg string)
	onClose   func(c *Client, err error)
}

// NewServer creates a new tcp server connection using the given net connection.
func NewServer(ip string, port int32) *Server {
	return &Server{
		ip:        ip,
		port:      port,
		onConnect: func(c *Client) {},
		onMessage: func(c *Client, msg string) {},
		onClose:   func(c *Client, err error) {},
		logger:    log.New(os.Stderr, "【F2F】", log.LstdFlags),
	}
}

func (s *Server) log(format string, v ...interface{}) {
	s.logger.Printf(format, v...)
}

// Start create a TCP server listener to accept client connection
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
			ExtraMap: map[string]string{},
		}
		go c.process()
	}
}

// OnConnect connect callbacks on a connection
func (s *Server) OnConnect(callback func(c *Client)) {
	s.onConnect = callback
}

// OnMessage receive callbacks on a connection
func (s *Server) OnMessage(callback func(c *Client, msg string)) {
	s.onMessage = callback
}

// OnClose close callbacks on a connection, it will be called before closing
func (s *Server) OnClose(callback func(c *Client, err error)) {
	s.onClose = callback
}

// Client defines parameters for accept an client
type Client struct {
	Conn     net.Conn
	server   *Server
	clientIP net.Addr
	ExtraMap map[string]string
}

// process client data
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

// Send send string message to client connection
func (c *Client) Send(msg string) error {
	return c.SendBytes([]byte(msg))
}

// SendLine send a line to client connection
func (c *Client) SendLine(line string) error {
	return c.SendBytes([]byte(line + "\r\n"))
}

// SendBytes send bytes message to client connection
func (c *Client) SendBytes(b []byte) error {
	data := Encode(string(b))
	_, err := c.Conn.Write(data)
	if err != nil {
		c.Conn.Close()
		c.server.onClose(c, err)
	}
	return err
}

// GetClientIP get client IP
func (c *Client) GetClientIP() net.Addr {
	return c.clientIP
}
