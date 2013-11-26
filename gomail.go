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

	var count int
	count = 1
	for {

		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		client := Client{
			Conn:     conn,
			Text:     textproto.NewConn(conn),
			didHelo:  false,
			ClientId: strconv.Itoa(count),
		}

		go hanndleConnection(client)
		count += 1
	}
}

func (c *Client) sendResponse(resp string) {
	resp = resp + " " + c.ClientId + " emel"
	id, err := c.Text.Cmd(resp)
	if err != nil {
		log.Println(err)
	}
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
}

func (c *Client) doHelo() {
	c.sendResponse("220 mx.valis.org SMTP gomail")
	c.didHelo = true
}

func (c *Client) Close() {
	c.sendResponse("221 2.0.0 closing connection")
	c.Text.Close()
	c.Conn.Close()
}

func (c *Client) sendHello() {
	c.sendResponse("250 mx.valis.org at your service")
}

func (c *Client) sendCommandNotRecognized() {
	c.sendResponse("500 unrecognized command")
}

func (c *Client) sendCommandNotImplemented() {
	c.sendResponse("502 5.5.1 Unimplemented comand")
}

func (c *Client) sendWillTryMyBest() {
	c.sendResponse("252 2.1.5 Send some mail, I'll try my best")
}

func (c *Client) sendOk() {
	c.sendResponse("250 2.1.5 OK")
}

func (c *Client) reset() {
	c.From = ""
	c.Rcpt = ""
	c.sendResponse("250 2.1.5 Flushed")
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
