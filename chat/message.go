package chat

import (
	"fmt"
)

type Message interface {
	String(user string) string
}

type SystemMessage struct {
	Message string
	Subject string
	NoteIfYou bool
	HideIfYou bool
}

type TextMessage struct {
	From string
	Text string
}

func (t SystemMessage) String(user string) string {
	isYou := user == t.Subject
	if isYou && t.HideIfYou {
		return ""
	}

	out := t.Message
	if isYou && t.NoteIfYou {
		return fmt.Sprintf("%s (** this is you)\n", out)
	} else {
		return fmt.Sprintf("%s\n", out)
	}
}

func (t TextMessage) String(user string) string {
	return fmt.Sprintf("%s: %s\n", t.From, t.Text)
}

