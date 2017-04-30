package main

import (
	"fmt"
	"net"
	"os"

	"github.com/yuichi1004/telnet-chat/chat/standalone"
)

const (
	CONN_PORT = "9399"
	CONN_TYPE = "tcp"
)

func main() {
	c := standalone.NewChat()

	host := os.Getenv("CHAT_HOST")
	port := os.Getenv("CHAT_PORT")
	if port == "" {
		port = CONN_PORT
	}

	l, err := net.Listen(CONN_TYPE, host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Listening on " + host + ":" + port)
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

