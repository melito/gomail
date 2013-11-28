package main

import (
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

type Server struct {
	Addr     string
	Name     string
	listener net.Listener
}

type Client struct {
	Conn     net.Conn
	Text     *textproto.Conn
	Rcpt     string
	From     string
	Data     []byte
	ClientId string
	server   *Server
	finished chan bool
}

func main() {
	server := newServer(3005)
	startServer(server)
}

func newServer(port int) Server {
	return Server{
		Addr: ":" + strconv.Itoa(port),
		Name: "localhost"}
}

func startServer(server Server) {
	server.Start()
}

func (server *Server) Start() {
	l, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
	defer l.Close()
	server.listener = l

	var count int
	count = 1
	for {

		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}

		client := Client{
			Conn:     conn,
			Text:     textproto.NewConn(conn),
			ClientId: strconv.Itoa(count),
			server:   server,
			finished: make(chan bool),
		}

		go hanndleConnection(client)

		count += 1
	}
}

func (server *Server) Stop() {
}

func (c *Client) sendResponse(resp string) {
	resp = resp + " - " + c.ClientId + " eMel"
	id, err := c.Text.Cmd(resp)
	if err != nil {
		log.Println(err)
	}
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
}

func (c *Client) doHelo() {
	c.sendResponse("220 " + c.server.Name + " SMTP")
}

func (c *Client) sendGoodbye() {
	c.sendResponse("221 2.0.0 closing connection")
	c.finished <- true
}

func (c *Client) Close() {
	log.Println("Closing connection", c.Conn.RemoteAddr())
	c.Text.Close()
	c.Conn.Close()
}

func (c *Client) sendHello() {
	c.sendResponse("250 " + c.server.Name + " at your service")
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

func (c *Client) sendGoAhead() {
	c.sendResponse("354 Enter message, ending with \".\" on a line by itself")
}

func (c *Client) getEmailData() {
	buf, err := c.Text.ReadDotBytes()
	if err != nil {
		log.Println(err)
	}
	c.Data = buf

	c.sendResponse("250 OK : Queued Message")
	log.Println("Queued message from", c.Conn.RemoteAddr())
}

func (c *Client) reset() {
	c.From = ""
	c.Rcpt = ""
	c.Data = nil
	c.sendResponse("250 2.1.5 Flushed")
}

func hanndleConnection(c Client) {
	defer c.Close()
	log.Println("Connection from", c.Conn.RemoteAddr())

	c.doHelo()

	for {
		select {
		case <-c.finished:
			break
		default:
			line, _ := c.Text.ReadLine()
			parseCommand(line, c)
		}
	}
}

func parseCommand(line string, c Client) {
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

	case "data":
		c.sendGoAhead()
		c.getEmailData()

	case "rset":
		c.reset()

	case "quit":
		c.sendGoodbye()

	default:
		c.sendCommandNotRecognized()
	}

}
