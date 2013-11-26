package main

import (
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

const (
	server_name = "mail.mel.io"
)

type Client struct {
	Conn    net.Conn
	Text    *textproto.Conn
	didHelo bool
	Rcpt    string
	From    string
	Data    string
}

func main() {
	startServer(3005)
}

func startServer(port int) {
	l, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		client := Client{Conn: conn, Text: textproto.NewConn(conn), didHelo: false}
		go hanndleConnection(client)
	}
}

func (c *Client) doHelo() {
	id, err := c.Text.Cmd("220 mx.valis.org SMTP gomail")
	if err != nil {
		log.Println(err)
	}
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)

	c.didHelo = true
}

func (c *Client) Close() {
	c.Text.Close()
	c.Conn.Close()
}

func hanndleConnection(c Client) {
	defer c.Close()
	log.Println("Connection from:", c.Conn.RemoteAddr())

	c.doHelo()

	for {
		line, _ := c.Text.ReadLine()
		finished := parseCommand(line, c)
		if finished {
			break
		}
	}
}

func parseCommand(line string, c Client) (finished bool) {
	pieces := strings.Split(line, " ")
	cmd := strings.ToLower(pieces[0])
	log.Println(cmd)

	switch cmd {
	case "helo":
		id, _ := c.Text.Cmd("250 mx.valis.org at your service")
		c.Text.StartResponse(id)
		defer c.Text.EndResponse(id)
		return false

	case "ehlo":
		id, _ := c.Text.Cmd("250 mx.valis.org at your service")
		c.Text.StartResponse(id)
		defer c.Text.EndResponse(id)
		return false
	case "quit":
		id, _ := c.Text.Cmd("221 2.0.0 closing connection")
		c.Text.StartResponse(id)
		defer c.Text.EndResponse(id)
		c.Close()
		return true
	}

	return false

}
