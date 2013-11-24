package main

import (
	"io"
	"log"
	"net"
	"strconv"
	"runtime"
)

const (
	server_name = "mail.mel.io"
)

var mainListener net.Listener

func main() {
	startServer(3005)
}

func startServer(port int) {
	log.Println(strconv.Itoa(port))
	mainListener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server started listening on port", strconv.Itoa(port))

	conn, err := mainListener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accepted connection from", conn.RemoteAddr())

	go handleConnection(conn)
}

func stopServer() {
	mainListener.Close()
}

func handleConnection(c net.Conn) {
	io.Copy(c, c)
	c.Close()
	log.Println("Connection closed", c.RemoteAddr())
}

func (c *net.Conn) server(){
	defer func(){
		if err := recover(); err != nil {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Println("smtpd: panic %v\n%s", err, buff)
		}
	}()
}
