package standalone

import (
	"github.com/yuichi1004/telnet-chat/chat"

	"fmt"
)

type Room struct {
	chat.Room
	pubsub *PubSub
}

type Chat struct {
	rooms map[string] *Room
}

type ChatParticipant struct {
	name string
	room *Room
	closer func()
}

// create new chat intance
func NewChat() *Chat {
	return &Chat{
		rooms: make(map[string] *Room),
	}
}

func (c *Chat) NewRoom(name string) error {
	_, ok := c.rooms[name]
	if ok {
		return fmt.Errorf("room already exists (name:%s)", name)
	}
	c.rooms[name] = &Room{
		chat.Room{
			Name: name, 
			Participants: make([]string, 0),
		},
		NewPubSub(),
	}
	return nil
}

func (c *Chat) GetRooms() ([]string, error) {
	names := make([]string, len(c.rooms))
	i := 0
	for _, r := range(c.rooms) {
		names[0] = r.Name
		i = i + 1
	}
	return names, nil
}

func (c *Chat) GetRoom(room string) (*chat.Room, error) {
	r, ok := c.rooms[room]
	if !ok {
		return nil, fmt.Errorf("room not found (name:%s)", room)
	}
	return &r.Room, nil
}

func (c *Chat) Join(room, user string) (chat.Participant, error) {
	r, ok := c.rooms[room]
	if !ok {
		return nil, fmt.Errorf("room not found (name:%s)", room)
	}

	p := &ChatParticipant {
		name: user,
		room: r,
	}
	return p, nil
}

func (p *ChatParticipant) Send(message string) error {
	return p.room.pubsub.Publish(message)
}

func (p *ChatParticipant) Subscribe()(chan string, error) {
	ch := make(chan string, 10)
	var err error
	p.closer, err = p.room.pubsub.Subscribe(ch)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (p *ChatParticipant) Leave() error {
	p.closer()
	return nil
}

