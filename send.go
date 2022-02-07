package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

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
	protocol network.Protocol
}

func send(ctx *cli.Context) (err error) {
	var sender = &Send{
		addr:     defaultAddr,
		protocol: network.NewDefaultProtocol(),
	}

	if ctx.String("ip") != defaultAddr {
		sender.addr, err = inputAddr()
		if err != nil {
			return
		}
	}
	sender.paths, err = inputFilePath()
	if err != nil {
		return
	}

	sender.conn, err = net.Dial("tcp", sender.addr)
	if err != nil {
		sender.logf("net dial err: %v", err)
		return
	}
	defer sender.conn.Close()

	ok, err := sender.ack()
	if err == nil && ok {
		for _, path := range sender.paths {
			sender.sendFile(path)
		}
	}
	return
}

func inputAddr() (addr string, err error) {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stdout, "Please enter the destination address\n\n(default: %s):\n", defaultAddr)
	addr, err = inputReader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("input params error: %v", err)
	}
	addr = strings.TrimSpace(addr)
	if len(addr) == 0 {
		addr = defaultAddr
	}
	return
}

func inputFilePath() (path []string, err error) {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stdout, "Please enter the file path to transfer(eg: "+"./test.txt"+"): \n")
	p, err := inputReader.ReadString('\n')
	if err != nil {
		err = fmt.Errorf("input params error: %v", err)
		return
	}

	p = strings.TrimSuffix(p, "\n")
	if p == "" {
		var f []byte
		f, err = ioutil.ReadFile("./ftf.conf")
		if err != nil {
			return
		}
		paths := strings.Split(string(f), ",")
		return paths, nil
	}
	path = []string{p}
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
