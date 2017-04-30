package redischat

import (
	"github.com/yuichi1004/telnet-chat/chat"
	"github.com/go-redis/redis"

	"fmt"
	"encoding/json"
)

type Instance struct {
	client *redis.Client
}

type ChatParticipant struct {
	client *redis.Client
	room string
	name string
	closed bool
	closer func() error
}


func NewInstance(host string) chat.Chat {
	return &Instance{
		client: redis.NewClient(&redis.Options{
			Addr:     host,
			Password: "",
			DB:       0,
		}),
	}
}

func (i *Instance) NewRoom(name string) error {
	exists, err := i.client.SIsMember("/rooms", name).Result()
	if err != nil {
		return fmt.Errorf("failed to get room info (err: %v)", err)
	}
	if exists {
		return fmt.Errorf("room already exists")
	}
	added, err := i.client.SAdd("/rooms", name).Result()
	if err != nil {
		return fmt.Errorf("failed to create room (err: %v)", err)
	}
	if added != 1 {
		return fmt.Errorf("room not created unexpectedly")
	}
	return nil
}
func (i *Instance) GetRooms() ([]string, error){
	rooms, err := i.client.SMembers("/rooms").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms (err: %v)", err)
	}
	return rooms, nil
}

func (i *Instance) GetRoom(room string) (*chat.Room, error){
	exists, err := i.client.SIsMember("/rooms", room).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get room info (err: %v)", err)
	}
	if !exists {
		return nil, fmt.Errorf("room already exists")
	}

	path := fmt.Sprintf("/rooms/%s/members", room)
	members, err := i.client.SMembers(path).Result()
	participants := make(map[string] chat.Participant, len(members))
	for _, v := range(members) {
		participants[v] = &ChatParticipant{name: v, room: room}
	}
	return &chat.Room{Name: room, Participants: participants}, nil
}

func (i *Instance) Join(room, user string) (chat.Participant, error) {
	path := fmt.Sprintf("/rooms/%s/members", room)
	if err := i.client.SAdd(path, user).Err(); err != nil {
		return nil, fmt.Errorf("join error for unexpected reason")
	}
	closer := func() error {
		return i.client.SRem(path, user).Err()
	}
	return &ChatParticipant{client: i.client, closer:closer, name: user, room: room}, nil
}

func (i *Instance) Connect(user string) error {
	exists, err := i.client.SIsMember("/users", user).Result()
	if err != nil {
		return fmt.Errorf("failed to user info (err: %v)", err)
	}
	if exists {
		return fmt.Errorf("user already exists")
	}
	added, err := i.client.SAdd("/users", user).Result()
	if err != nil {
		return fmt.Errorf("failed to create room (err: %v)", err)
	}
	if added != 1 {
		return fmt.Errorf("room not created unexpectedly")
	}
	return nil
}

func (i *Instance) Disconnect(user string) error {
	return i.client.SRem("/users", user).Err()
}

func (c *ChatParticipant) Send(message chat.Message) error {
	data, err := json.Marshal(&message)
	if err != nil {
		return err
	}
	subpath := fmt.Sprintf("/rooms/%s/subscribers", c.room)
	return c.client.Publish(subpath, string(data)).Err()
}

func (c *ChatParticipant) Subscribe(ch chan chat.Message) error{
	go func() {
		subpath := fmt.Sprintf("/rooms/%s/subscribers", c.room)
		pubsub := c.client.Subscribe(subpath)
		defer pubsub.Close()

		str, _ := pubsub.ReceiveMessage()
		for !c.closed {
			msg, err := chat.FromJSON(str.Payload)
			if err != nil {
				fmt.Printf("error: %+v\n", err)
			}
			if msg != nil {
				ch <- msg
			}
			str, _ = pubsub.ReceiveMessage()
		}
	}()
	return nil
}

func (c *ChatParticipant) Leave() error {
	c.closed = true
	return c.closer()
}

func (c *ChatParticipant) Name() string{
	return c.name
}

func (c *ChatParticipant) Room() string{
	return c.room
}

