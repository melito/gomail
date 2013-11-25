package main

import (
	"log"
	"net"
	"runtime"
)

const (
	server_name = "mail.mel.io"
)

type Server struct {
	Addr string
}

type conn struct {
	remoteAddr string
	server     *Server
	rwc        net.Conn
}

func ListenAndServe(addr string) error {
	server := &Server{Addr: addr}
	return server.ListenAndServe()
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":stmp"
	}
	l, e := net.Listen("tcp", addr)
	if e != nil {
		return e
	}
	return srv.Serve(l)
}

func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()

	go func() {
		for {
			rw, e := l.Accept()
			if e != nil {
				return e
			}
			c, err := srv.newConn(rw)
			if err != nil {
				log.Println(err)
				continue
			}
			go c.serve()
		}
	}()
}

func (srv *Server) newConn(rwc net.Conn) (c *conn, err error) {
	c = new(conn)
	c.remoteAddr = rwc.RemoteAddr().String()
	c.server = srv
	c.rwc = rwc
	return c, nil
}

func (c *conn) serve() {
	defer func() {
		if err := recover(); err != nil {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Println(err, buf)
		}
	}()
}
