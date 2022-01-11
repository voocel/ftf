package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Send struct {
	addr   string
	paths  []string
	conn   net.Conn
	logger *log.Logger
}

func send(c *cli.Context) (err error) {
	var sender = Send{
		addr: defaultAddr,
	}

	if c.String("ip") != defaultAddr {
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
		sender.logf("net dial err: ", err)
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
		f, err = os.ReadFile("./f2f.conf")
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
		s.logf("os stat err: ", err)
		return
	}

	_, err = s.conn.Write([]byte(info.Name()))
	if err != nil {
		s.logf("conn.Write info.Name err: ", err)
		return
	}

	var n int
	buf := make([]byte, 1024)
	s.log("waiting ack···")
	n, err = s.conn.Read(buf)
	if err != nil {
		s.logf("conn read ok err: ", err)
		return
	}
	if "ok" == string(buf[:n]) {
		s.log("ack success")
		return true, nil
	}
	return
}

func (s *Send) sendFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		s.logf("os open err: ", err)
		return
	}
	defer file.Close()
	buf := make([]byte, 1024*4)

	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				s.log("files send success")
			} else {
				s.logf("file read err: ", err)
			}
			return
		}
		if n == 0 {
			fmt.Println("files send success!!!")
			break
		}
		s.conn.Write(buf[:n])
	}
}

func (s *Send) log(v ...interface{}) {
	if s.logger == nil {
		s.logger = log.New(os.Stderr, "【F2F】", log.LstdFlags)
	}
	s.logger.Print(v...)
}

func (s *Send) logf(format string, v ...interface{}) {
	if s.logger == nil {
		s.logger = log.New(os.Stderr, "【F2F】", log.LstdFlags)
	}
	s.logger.Printf(format, v...)
}
