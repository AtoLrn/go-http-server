package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go printHttpRequest(conn)
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
	headers interface{}
	body    interface{}
}

func createHttpObject(req string) Request {
	lines := strings.Split(req, "\n")
	method := strings.Split(lines[0], " ")[0]
	return Request{method: method}
}

func printHttpRequest(conn net.Conn) {
	defer conn.Close()
	for {
		messageByte, err := GetWholeBuffer(bufio.NewReader(conn))
		if err != nil {
			fmt.Println("Connection closed.")
			return
		}
		req := createHttpObject(messageByte)

		fmt.Println(req)
		conn.Write([]byte("Hello"))
		conn.Close()
	}
}
