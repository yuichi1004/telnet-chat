package chat

type Room struct {
	Name string
	Participants map[string] Participant
}

type Chat interface {
	// Create new chat room
	NewRoom(name string) error

	// Get list of rooms
	GetRooms() ([]string, error)

	// Get specific room
	GetRoom(room string) (*Room, error)
	
	// Join to a chat room as a participant
	Join(room, user string) (Participant, error)

	// connect a user
	Connect(user string) error

	// diconnect a user
	Disconnect(user string) error
}


type Participant interface {
	// Send a message as a participant
	Send(message string) error

	// Subscribe messages on the chat room
	Subscribe(ch chan string) error

	// Leave the chat room
	Leave() error

	// Get name of participant
	Name() string

	// Get room name
	Room() string
}
