package main

import (
	"github.com/yuichi1004/telnet-chat/chat"

	"fmt"
	"io"
	"bufio"
	"strings"
)
const (
)

type ChatHandler struct {
	chat chat.Chat
	participant chat.Participant
	writer io.ReadWriteCloser
	subscriber chan string

	name string
	connected bool
}

func NewChatHandler(chat chat.Chat, writer io.ReadWriteCloser) *ChatHandler {
	c:=  &ChatHandler {
		chat: chat,
		writer: writer,
		subscriber: make(chan string),
		connected: true,
	}
	c.printf("Welcome to the XYZ chat server\n")
	c.printf("Loigin name?\n")
	return c
}

func (c *ChatHandler) doHandle() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover on doHandle()")
		}
	}()
	defer c.writer.Close()

	line := make(chan string)
	go func() {
		reader := bufio.NewReader(c.writer)
		for c.connected {
			buf, _, err := reader.ReadLine()
			switch err {
			case nil:
			case io.EOF:
				c.disconnect()
				fmt.Println("Connection closed")
				return
			default:
				fmt.Println("Error reading:", err.Error())
				return
			}
			line <- string(buf)
		}
	}()

	for c.connected {
		select {
		case msg:=<-c.subscriber:
			if strings.HasSuffix(msg, c.name) && strings.HasPrefix(msg, " *") {
				c.printf("%s (** this is you)\n", msg)
			} else {
				c.printf("%s\n", msg)
			}
		case msg:=<-line:
			c.HandleLine(msg)
		}
	}
}

func (c *ChatHandler) HandleLine(line string) error {
	defer func() {
		if r := recover(); r != nil {
			c.quit()
		}
	}()

	if c.name == "" {
		if err := c.chat.Connect(line); err != nil {
			c.printf("Sorry, name taken.\nLogin name?\n")
			return nil
		}
		c.name = line
		c.printf("Welcome %s!\n", c.name)
		return nil
	}


	var err error
	if strings.HasPrefix(line, "/") {
		args := strings.Split(line, " ")
		if len(args) == 0 {
			return nil
		}
		switch args[0] {
		case "/rooms":
			err = c.rooms()
		case "/join":
			if len(args) < 2 {
				err = fmt.Errorf("please type room name")
			} else {
				err = c.join(args[1])
			}
		case "/leave":
			err = c.leave()
		case "/quit":
			err = c.quit()
		case "/help":
			err = c.help()
		default:
			err = fmt.Errorf("unknown command: %s", args[0])
		}
	} else {
		err = c.send(line)
	}

	if err != nil {
		c.printf("error: %v\n", err)
	}

	return nil
}

func (c *ChatHandler) send(msg string) error {
	if c.participant == nil {
		return fmt.Errorf("you need to join the room first")
	}
	text := fmt.Sprintf("%s: %s", c.name, msg)
	return c.participant.Send(text)
}

func (c *ChatHandler) rooms() error {
	rooms, err := c.chat.GetRooms()
	if err != nil {
		return fmt.Errorf("failed to get rooms")
	}
	for _, r := range rooms {
		room, _ := c.chat.GetRoom(r)
		c.printf("* %s (%d)\n", r, len(room.Participants))
	}
	c.printf("end of list\n")

	return nil
}

func (c *ChatHandler) join(room string) error {
	if c.participant != nil {
		return fmt.Errorf("your are already on %s. /leave the room first", room)
	}

	r, err := c.chat.GetRoom(room)
	if err != nil {
		c.chat.NewRoom(room)
		r, err = c.chat.GetRoom(room)
		if err != nil {
			return fmt.Errorf("failed to create new room")
		}
	}
	c.participant, err = c.chat.Join(room, c.name)
	if err != nil {
		return fmt.Errorf("failed to subscribe the room")
	}
	c.subscriber, err = c.participant.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe the room")
	}

	c.printf("entering room: %s\n", room)
	for _, p := range(r.Participants) {
		note := ""
		if p.Name() == c.name {
			note = " (** this is you)"
		}
		c.printf(" * %s%s\n", p.Name(), note)
	}
	c.printf("end of list\n")

	msg := fmt.Sprintf(" * new user joined  %s: %s", c.participant.Room(), c.name)
	c.participant.Send(msg)

	return nil
}

func (c *ChatHandler) leave() error {
	msg := fmt.Sprintf(" * user has left %s: %s", c.participant.Room(), c.name)
	c.participant.Send(msg)
	err := c.participant.Leave()
	c.participant = nil
	return err
}

func (c *ChatHandler) disconnect() error {
	if c.name != "" {
		c.chat.Disconnect(c.name)
	}
	return nil
}

func (c *ChatHandler) quit() error{
	c.printf("BYE\n")
	c.connected = false
	c.disconnect()
	return nil
}

func (c *ChatHandler) help() error {
	c.printf(" * commands:\n")
	c.printf("    - /join [room] - join chat room\n")
	c.printf("    - /rooms       - list chat rooms\n")
	c.printf("    - /leave       - leave the current chat room\n")
	c.printf("    - /quit        - quit this chat server\n")
	return nil
}

func (c *ChatHandler) printf(format string, args... interface{}) {
	text := fmt.Sprintf(format, args...)
	c.writer.Write([]byte(text))
}

