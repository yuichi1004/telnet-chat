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
	users map[string] bool
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
		users: make(map[string] bool),
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
			Participants: make(map[string]chat.Participant, 0),
		},
		NewPubSub(),
	}
	return nil
}

func (c *Chat) GetRooms() ([]string, error) {
	names := make([]string, len(c.rooms))
	i := 0
	for _, r := range(c.rooms) {
		names[i] = r.Name
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
	r.Participants[user] = p
	return p, nil
}

func (c *Chat) Connect(user string) error {
	if _, ok := c.users[user]; ok {
		return fmt.Errorf("user %s already exists", user)
	}
	c.users[user] = true
	return nil
}

func (c *Chat) Disconnect(user string) error {
	if _, ok := c.users[user]; ok {
		delete(c.users, user)
	}
	return nil
}

func (p *ChatParticipant) Send(message chat.Message) error {
	return p.room.pubsub.Publish(message)
}

func (p *ChatParticipant) Subscribe(ch chan chat.Message)(error) {
	var err error
	p.closer, err = p.room.pubsub.Subscribe(ch)
	if err != nil {
		return err
	}
	return nil
}

func (p *ChatParticipant) Leave() error {
	p.closer()
	delete(p.room.Participants, p.name)
	return nil
}

func (p *ChatParticipant) Name() string {
	return p.name
}

func (p *ChatParticipant) Room() string {
	return p.room.Name
}
