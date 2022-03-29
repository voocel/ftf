package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/voocel/ftf/network"
)

const (
	DefaultRecvBufSize = 4 << 10 // 4k
)

type Send struct {
	addr     string
	paths    []string
	conn     net.Conn
	logger   *log.Logger
	folder   bool
	protocol network.Protocol
}

func NewSend(addr string) *Send {
	return &Send{
		addr:     addr,
		protocol: network.NewDefaultProtocol(),
	}
}

func sender(ctx *cli.Context) (err error) {
	if err := runApp(ctx); err != nil {
		return cli.Exit(err, 1)
	}
	return
}

func (s *Send) ack() (ok bool, err error) {
	info, err := os.Stat(s.paths[0])
	if err != nil {
		s.logf("os stat err: %v", err)
		return
	}

	data := network.NewMessage(network.Ack, []byte(info.Name()))
	msg, err := s.protocol.Pack(data)
	if err != nil {
		s.logf("ack packet err: %v", err)
		return
	}
	_, err = s.conn.Write(msg)
	if err != nil {
		s.logf("send ack err: %v", err)
		return
	}

	s.log("waiting for ack···")
	reader := bufio.NewReader(s.conn)
	res, err := s.protocol.Unpack(reader)
	if err != nil {
		s.logf("receive ack err: %v", err)
		return
	}
	if res.GetCmd() == network.Ack {
		s.log("ack success")
		return true, nil
	}
	return
}

func (s *Send) sendFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		s.logf("os open err: %v", err)
		return
	}
	defer file.Close()

	buf := make([]byte, DefaultRecvBufSize)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				s.log("files send success")
			} else {
				s.logf("file read err: %v", err)
			}
			return
		}
		if n == 0 {
			s.log("files send success!!!")
			break
		}
		msg := network.NewMessage(network.Single, buf[:n])
		b, err := s.protocol.Pack(msg)
		if err != nil {
			s.logf("send file message err: %v", err)
			return
		}
		s.conn.Write(b)
	}
}

func (s *Send) addPath(path ...string) {
	s.paths = append(s.paths, path...)
}

func (s *Send) log(v ...interface{}) {
	if s.logger == nil {
		s.logger = log.New(os.Stderr, "【FTF】", log.LstdFlags)
	}
	s.logger.Print(v...)
}

func (s *Send) logf(format string, v ...interface{}) {
	if s.logger == nil {
		s.logger = log.New(os.Stderr, "【FTF】", log.LstdFlags)
	}
	s.logger.Printf(format, v...)
}
