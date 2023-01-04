package lib

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type HTTPServer struct {
	Port     int
	Router   Router
	Listener net.Listener
}

func (h HTTPServer) Run() {
	ln, err := net.Listen("tcp", ":"+fmt.Sprint(h.Port))
	if err != nil {
		fmt.Println(err)
		return
	}

	h.Listener = ln
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go h.printHttpRequest(conn)
	}
}

func GetWholeBuffer(buff *bufio.Reader) (string, error) {
	content := ""
	for {
		line, _, err := buff.ReadLine()
		if err != nil {
			return content, errors.New("a")
		}
		content += (string(line) + "\n")

		if buff.Buffered() == 0 {
			break
		}
	}
	return content, nil
}

type Request struct {
	method  string
	headers map[string]string
	body    string
	Path    string
	Conn    net.Conn
}

func createHttpObject(req string, conn net.Conn) Request {
	step := 0
	lines := strings.Split(req, "\n")
	firstLine := strings.Split(lines[0], " ")
	method := firstLine[0]
	path := firstLine[1]
	headers := make(map[string]string)
	body := ""
	for i, s := range lines {
		if i == 0 {
			continue
		}

		if s == "" {
			step = 1
		}
		if step == 0 {
			header := strings.Split(s, ": ")
			headers[header[0]] = header[1]
		}
		if step == 1 {
			body += s
		}

	}
	return Request{method: method, headers: headers, body: body, Path: path, Conn: conn}
}

func (h HTTPServer) AddRoute(r string, cb func(req Request, res Response)) {
	h.Router.Routes[r] = cb
}

func (h HTTPServer) printHttpRequest(conn net.Conn) {
	defer conn.Close()
	for {
		messageByte, err := GetWholeBuffer(bufio.NewReader(conn))
		if err != nil {
			fmt.Println("Connection closed.")
			return
		}
		req := createHttpObject(messageByte, conn)
		res := Response{Conn: conn, Request: req}
		if h.Router.Execute(req, res) == false {
			res.NotFound()
		}
		conn.Close()
	}
}

func CreateServer(port int) HTTPServer {
	server := HTTPServer{}

	server.Port = port
	server.Router = GetRouter()
	return server
}
