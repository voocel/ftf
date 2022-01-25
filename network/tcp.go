package network

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// Server defines parameters for running an TCP network
type Server struct {
	ip        string
	port      int32
	logger    *log.Logger
	tlsConf   *tls.Config
	onConnect func(c *Client)
	onMessage func(c *Client, msg *Message)
	onClose   func(c *Client, err error)
}

// NewServer creates a new tcp network connection using the given net connection.
func NewServer(ip string, port int32, opts ...Options) *Server {
	serv := &Server{
		ip:        ip,
		port:      port,
		onConnect: func(c *Client) {},
		onMessage: func(c *Client, msg *Message) {},
		onClose:   func(c *Client, err error) {},
		logger:    log.New(os.Stderr, "【FTF】", log.LstdFlags),
	}
	for _, opt := range opts {
		opt(serv)
	}
	return serv
}

func (s *Server) log(format string, v ...interface{}) {
	s.logger.Printf(format, v...)
}

// Start create a TCP network listener to accept client connection
func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	s.log("network start success! %s:%d", s.ip, s.port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
			protocol: NewDefaultProtocol(),
			ExtraMap: map[string]string{},
		}
		go c.process(ctx)
	}
}

// OnConnect connect callbacks on a connection
func (s *Server) OnConnect(callback func(c *Client)) {
	s.onConnect = callback
}

// OnMessage receive callbacks on a connection
func (s *Server) OnMessage(callback func(c *Client, msg *Message)) {
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
	protocol Protocol
	timer    *time.Timer
	timeout  time.Duration
	interval time.Duration
	sendCh   chan *Message
	msgCh    chan *Message
	errDone  chan error
	ExtraMap map[string]string
}

// process client connection
func (c *Client) process(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		c.Conn.Close()
	}()

	go c.readLoop(ctx)
	go c.writeLoop(ctx)

	c.server.onConnect(c)
	for {
		select {
		case err := <-c.errDone:
			c.server.onClose(c, err)
			return
		case msg := <-c.msgCh:
			c.server.onMessage(c, msg)
		}
	}
}

// readLoop read goroutine
func (c *Client) readLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			reader := bufio.NewReader(c.Conn)
			msg, err := c.protocol.Unpack(reader)
			if err != nil {
				c.errDone <- err
				continue
			}
			// todo skip ping
			c.msgCh <- msg
		}
	}
}

// writeLoop write goroutine
func (c *Client) writeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.sendCh:
			if err := c.Send(msg); err != nil {
				log.Printf("send message err: %v", err)
			}
		case <-c.timer.C:
			msg := NewMessage(Heartbeat, []byte("ping"))
			c.SendMessage(msg)
			if c.interval > 0 {
				c.timer.Reset(c.interval)
			}
		}
	}
}

// Send send string message to client connection
func (c *Client) Send(msg *Message) error {
	m, err := c.protocol.Pack(msg)
	if err != nil {
		return err
	}
	return c.SendBytes(m)
}

// SendLine send a line to client connection
func (c *Client) SendLine(line string) error {
	return c.SendBytes([]byte(line + "\r\n"))
}

// SendBytes send bytes message to client connection
func (c *Client) SendBytes(b []byte) error {
	_, err := c.Conn.Write(b)
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

// SendMessage send message into channel
func (c *Client) SendMessage(msg *Message) {
	c.sendCh <- msg
}
