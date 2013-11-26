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
	Conn     net.Conn
	Text     *textproto.Conn
	didHelo  bool
	Rcpt     string
	From     string
	Data     string
	ClientId string
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

		client := Client{
			Conn:    conn,
			Text:    textproto.NewConn(conn),
			didHelo: false,
		}

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
	id, _ := c.Text.Cmd("221 2.0.0 closing connection")
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
	c.Text.Close()
	c.Conn.Close()
}

func (c *Client) sendHello() {
	id, _ := c.Text.Cmd("250 mx.valis.org at your service")
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
}

func (c *Client) sendCommandNotRecognized() {
	id, _ := c.Text.Cmd("500 unrecognized command")
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
}

func (c *Client) sendCommandNotImplemented() {
	id, _ := c.Text.Cmd("502 5.5.1 Unimplemented comand")
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
}

func (c *Client) sendWillTryMyBest() {
	id, _ := c.Text.Cmd("252 2.1.5 Send some mail, I'll try my best")
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
}

func (c *Client) sendOk() {
	id, _ := c.Text.Cmd("250 2.1.5 OK")
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
}

func (c *Client) reset() {
	c.From = ""
	c.Rcpt = ""
	id, _ := c.Text.Cmd("250 2.1.5 Flushed")
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
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
	//log.Println(line)

	switch cmd {
	case "helo":
		c.sendHello()

	case "ehlo":
		c.sendHello()

	case "expn":
		c.sendCommandNotImplemented()

	case "vrfy":
		c.sendWillTryMyBest()

	case "rcpt":
		c.sendOk()

	case "mail":
		c.sendOk()

	case "rset":
		c.reset()

	case "quit":
		c.Close()
		return true

	default:
		c.sendCommandNotRecognized()
	}

	return false

}
