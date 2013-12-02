package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/textproto"
	"path/filepath"
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
	Rcpt     []*mail.Address
	From     *mail.Address
	Data     []byte
	ClientId string
	server   *Server
	finished chan bool
}

var (
	default_server_name = "localhost"
	default_listen_port = 25
)

func init() {
	var user, pass string
	flag.IntVar(&default_listen_port, "port", 25, "default port to listen on")
	flag.StringVar(&default_server_name, "name", "localhost", "default server name used in banner")
	flag.StringVar(&user, "u", "user@example.com", "Add a user")
	flag.StringVar(&pass, "x", "", "Set a password (used with the -u flag when adding a user")
	flag.Parse()

	if pass != "" {
		if user == "" {
			log.Fatal("You set a password, but didn't set a user")
		}

		err := addUser(user, pass)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Added user", user)
		return
	}

}

func main() {
	server := newServer(default_listen_port)
	startServer(server)
}

func newServer(port int) Server {
	return Server{
		Addr: ":" + strconv.Itoa(port),
		Name: default_server_name}
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

	c.deliverEmail()

	c.sendResponse("250 OK : Queued Message")
	log.Println("Queued message from", c.Conn.RemoteAddr())
}

func (c *Client) deliverEmail() {

	for _, userAddr := range c.Rcpt {
		path := filepath.Join(pathToMailDirForEmail(userAddr.Address), "new")
		filename := filepath.Join(path, createUniqueFileName())
		err := ioutil.WriteFile(filename, c.Data, 0700)
		if err != nil {
			log.Println(err)
		}

	}

}

func (c *Client) reset() {
	c.From = nil
	c.Rcpt = nil
	c.Data = nil
	c.sendResponse("250 2.1.5 Flushed")
}

func (c *Client) parseRcptData(line string) {
	c.Rcpt, _ = mail.ParseAddressList(line)
}

func (c *Client) parseMailFromData(line string) {
	c.From, _ = mail.ParseAddress(line)
}

func getDataFromLineAfterColon(line string) string {
	pieces := strings.Split(line, ":")
	addressStr := strings.Join(pieces[1:], "")
	return addressStr
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
		c.parseRcptData(getDataFromLineAfterColon(line))
		c.sendOk()

	case "mail":
		c.parseMailFromData(getDataFromLineAfterColon(line))
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
