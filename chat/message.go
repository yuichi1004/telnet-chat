package chat

import (
	"fmt"
	"encoding/json"
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

func (t SystemMessage) MarshalJSON() ([]byte, error) {
	type alias SystemMessage
	v := struct {
		alias
		Type string
	}{
		alias(t),
		"system",
	}
	return json.Marshal(&v)
}

func (t TextMessage) String(user string) string {
	return fmt.Sprintf("%s: %s\n", t.From, t.Text)
}

func (t TextMessage) MarshalJSON() ([]byte, error) {
	type alias TextMessage
	v := struct {
		alias
		Type string
	}{
		alias(t),
		"text",
	}
	return json.Marshal(&v)
}

func FromJSON(js string) (Message, error) {
	v := struct {
		Type string
	}{}
	data := []byte(js)
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	switch v.Type {
	case "system":
		v := SystemMessage{}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return v, nil
	case "text":
		v := TextMessage{}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown message type: %s", v.Type)
	}
}
