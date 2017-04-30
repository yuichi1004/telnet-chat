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

	chat := NewChat()
	err := chat.NewRoom(room)
	if err != nil {
		t.Errorf("failed to create room (err:%v)", err)
	}

	p1, err := chat.Join(room, user1)
	if err != nil {
		t.Errorf("failed to join the room (err:%v)", err)
	}
	p2, err := chat.Join(room, user2)
	if err != nil {
		t.Errorf("failed to join the room (err:%v)", err)
	}

	ch1 := make(chan string, 10)
	ch2 := make(chan string, 10)
	p1.Subscribe(ch1)
	p2.Subscribe(ch2)

	check := func(expects string, ch chan string) error {
		select {
		case got:= <-ch:
			if got != expects {
				return fmt.Errorf("unexpected message (expects:%s, got:%s)", expects, got)
			}
			return nil
		default:
			return fmt.Errorf("failed to get message")
		}
	}

	p1.Send("john: hello")
	if err := check("john: hello", ch1); err != nil {
		t.Errorf("%v", err)
	}
	if err := check("john: hello", ch2); err != nil {
		t.Errorf("%v", err)
	}

	p2.Send("mike: hello")
	if err := check("mike: hello", ch1); err != nil {
		t.Errorf("%v", err)
	}
	if err := check("mike: hello", ch2); err != nil {
		t.Errorf("%v", err)
	}

	p1.Leave()

	p2.Send("mike: hello")
	if err := check("mike: hello", ch1); err == nil {
		t.Errorf("expected to be failed to get message but not")
	}
	if err := check("mike: hello", ch2); err != nil {
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


