package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net"
	"os"
)

func receive(c *cli.Context) error {
	var addr = "0.0.0.0:1234"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Println("waiting for connect···")
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	fileName := string(buf[:n])

	n, err = conn.Write([]byte("ok"))
	if err != nil {
		return
	}

	recvFile(fileName, conn)
}

func recvFile(fileName string, conn net.Conn) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("create file err: ", err)
		return
	}
	defer file.Close()

	buf := make([]byte, 1024*4)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("file receive success")
			} else {
				fmt.Println("conn.Read err =", err)
			}
			return
		}

		n, err = file.Write(buf[:n])
		if err != nil {
			log.Fatalln("file write err: ", err)
			return
		}
	}
}
