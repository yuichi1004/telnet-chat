package main

import (
	"fmt"
	"net"
	"os"

	"github.com/yuichi1004/telnet-chat/chat/standalone"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "9399"
	CONN_TYPE = "tcp"
)

func main() {
	c := standalone.NewChat()

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
		h := NewChatHandler(c, conn)
		go h.doHandle()
	}
}

