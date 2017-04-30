package standalone

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yuichi1004/telnet-chat/chat"
)

func TestChat(t *testing.T) {
	room := "hobby"
	user1 := "john"
	user2 := "mike"

	chatInstance := NewChat()
	err := chatInstance.NewRoom(room)
	if err != nil {
		t.Errorf("failed to create room (err:%v)", err)
	}

	p1, err := chatInstance.Join(room, user1)
	if err != nil {
		t.Errorf("failed to join the room (err:%v)", err)
	}
	p2, err := chatInstance.Join(room, user2)
	if err != nil {
		t.Errorf("failed to join the room (err:%v)", err)
	}

	ch1 := make(chan chat.Message, 10)
	ch2 := make(chan chat.Message, 10)
	p1.Subscribe(ch1)
	p2.Subscribe(ch2)

	check := func(expects chat.Message, ch chan chat.Message) error {
		select {
		case got:= <-ch:
			if !reflect.DeepEqual(got, expects) {
				return fmt.Errorf("unexpected message (expects:%s, got:%s)", expects, got)
			}
			return nil
		default:
			return fmt.Errorf("failed to get message")
		}
	}

	msg1 := chat.TextMessage{"john", "hello"}
	p1.Send(msg1)
	if err := check(msg1, ch1); err != nil {
		t.Errorf("%v", err)
	}
	if err := check(msg1, ch2); err != nil {
		t.Errorf("%v", err)
	}

	msg2 := chat.TextMessage{"mike", "hello"}
	p2.Send(msg2)
	if err := check(msg2, ch1); err != nil {
		t.Errorf("%v", err)
	}
	if err := check(msg2, ch2); err != nil {
		t.Errorf("%v", err)
	}

	p1.Leave()

	p2.Send(msg2)
	if err := check(msg2, ch1); err == nil {
		t.Errorf("expected to be failed to get message but not")
	}
	if err := check(msg2, ch2); err != nil {
		t.Errorf("%v", err)
	}
}

func TestChat_GetRooms(t *testing.T) {
	type fields struct {
		rooms map[string]*Room
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chat{
				rooms: tt.fields.rooms,
			}
			got, err := c.GetRooms()
			if (err != nil) != tt.wantErr {
				t.Errorf("Chat.GetRooms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chat.GetRooms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChat_GetRoom(t *testing.T) {
	type fields struct {
		rooms map[string]*Room
	}
	type args struct {
		room string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *chat.Room
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chat{
				rooms: tt.fields.rooms,
			}
			got, err := c.GetRoom(tt.args.room)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chat.GetRoom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chat.GetRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}


