package main

import (
	"fmt"
	"net"
	"os"
	"io"
	"bufio"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "9399"
	CONN_TYPE = "tcp"
)

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		buf, _, err := reader.ReadLine()
		switch err {
		case nil:
		case io.EOF:
			fmt.Println("Connection closed")
			return
		default:
			fmt.Println("Error reading:", err.Error())
			return
		}
		conn.Write([]byte(fmt.Sprintf("echo: %s\n", buf)))
	}
}
