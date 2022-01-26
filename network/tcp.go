package network

import (
	"bufio"
	"context"
	"crypto/tls"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// Server defines parameters for running an TCP network
type Server struct {
	addr      string
	logger    *log.Logger
	tlsConf   *tls.Config
	exitCh    chan struct{}
	onConnect func(c *Client)
	onMessage func(c *Client, msg *Message)
	onClose   func(c *Client, err error)
}

// NewServer creates a new tcp network connection using the given net connection.
func NewServer(addr string, opts ...Options) *Server {
	serv := &Server{
		addr:      addr,
		exitCh:    make(chan struct{}),
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
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	s.log("network start success! ", s.addr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case <-s.exitCh:
			return
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			s.log("accept connection err: %v", err)
			continue
		}
		c := Client{
			srv:      s,
			conn:     conn,
			timer:    time.NewTimer(5 * time.Second),
			clientIP: conn.RemoteAddr(),
			protocol: NewDefaultProtocol(),
			msgCh:    make(chan *Message, 1024),
			sendCh:   make(chan *Message, 1024),
			errDone:  make(chan error),
			extraMap: map[string]string{},
		}
		go c.process(ctx)
	}
}

// Stop Close server
func (s *Server) Stop() {
	close(s.exitCh)
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
	srv      *Server
	conn     net.Conn
	clientIP net.Addr
	protocol Protocol
	timer    *time.Timer
	timeout  time.Duration
	interval time.Duration
	sendCh   chan *Message
	msgCh    chan *Message
	errDone  chan error
	extraMap map[string]string
	sync.RWMutex
}

// process client connection
func (c *Client) process(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		c.conn.Close()
	}()

	go c.readLoop(ctx)
	go c.writeLoop(ctx)

	c.srv.onConnect(c)
	for {
		select {
		case <-c.srv.exitCh:
			return
		case err := <-c.errDone:
			c.srv.onClose(c, err)
			return
		case msg := <-c.msgCh:
			c.srv.onMessage(c, msg)
		}
	}
}

// readLoop read goroutine
func (c *Client) readLoop(ctx context.Context) {
	for {
		select {
		case <-c.srv.exitCh:
			return
		case <-ctx.Done():
			return
		default:
			reader := bufio.NewReader(c.conn)
			msg, err := c.protocol.Unpack(reader)
			if err != nil {
				c.errDone <- err
				continue
			}
			c.msgCh <- msg
		}
	}
}

// writeLoop write goroutine
func (c *Client) writeLoop(ctx context.Context) {
	for {
		select {
		case <-c.srv.exitCh:
			return
		case <-ctx.Done():
			return
		case msg := <-c.sendCh:
			if err := c.writeMessage(msg); err != nil {
				log.Printf("send message err: %v", err)
			}
		case <-c.timer.C:
			c.SendBytes(Heartbeat, []byte("ping"))
			if c.interval > 0 {
				c.timer.Reset(c.interval)
			}
		}
	}
}

// write Message to client connection
func (c *Client) writeMessage(msg *Message) error {
	m, err := c.protocol.Pack(msg)
	if err != nil {
		return err
	}
	return c.writeBytes(m)
}

// writeBytes send bytes message to client connection
func (c *Client) writeBytes(b []byte) error {
	_, err := c.conn.Write(b)
	if err != nil {
		c.conn.Close()
		c.srv.onClose(c, err)
	}
	return err
}

// GetRawConn get the raw net.Conn from the client connection
func (c *Client) GetRawConn() net.Conn {
	return c.conn
}

// SetExtraMap set the extra data
func (c *Client) SetExtraMap(k, v string) {
	c.Lock()
	defer c.Unlock()
	c.extraMap[k] = v
}

// GetExtraMap get the extra data
func (c *Client) GetExtraMap(k string) string {
	c.RLock()
	defer c.RUnlock()
	return c.extraMap[k]
}

// GetClientIP get client IP
func (c *Client) GetClientIP() net.Addr {
	return c.clientIP
}

// SendMessage send message into channel
func (c *Client) SendMessage(msg *Message) {
	c.sendCh <- msg
}

// SendBytes send bytes
func (c *Client) SendBytes(cmd CMD, b []byte) {
	msg := NewMessage(cmd, b)
	c.SendMessage(msg)
}
