package standalone

import(
	"github.com/yuichi1004/telnet-chat/chat"

	"sync"
)

type PubSub struct {
	subs map[int]chan chat.Message
	m sync.Mutex
	serial int
}

func NewPubSub() (*PubSub) {
	return &PubSub{
		subs: make(map[int]chan chat.Message),
	}
}

func (p *PubSub) Subscribe(c chan chat.Message) (func(), error) {
	p.m.Lock()
	defer p.m.Unlock()

	id := p.serial
	p.serial = p.serial + 1
	p.subs[id] = c

	unsub := func() {
		p.m.Lock()
		defer p.m.Unlock()
		delete(p.subs, id)
	}
	return unsub, nil
}

func (p *PubSub) Publish(message chat.Message) error {
	p.m.Lock()
	defer p.m.Unlock()
	for _, c := range p.subs {
		c <- message
	}

	return nil
}

