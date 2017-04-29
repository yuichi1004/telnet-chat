package standalone

import(
	"sync"
)

type PubSub struct {
	subs map[int]chan string
	m sync.Mutex
	serial int
}

func NewPubSub() (*PubSub) {
	return &PubSub{
		subs: make(map[int]chan string),
	}
}

func (p *PubSub) Subscribe(c chan string) (func(), error) {
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

func (p *PubSub) Publish(message string) error {
	p.m.Lock()
	defer p.m.Unlock()
	for _, c := range p.subs {
		c <- message
	}

	return nil
}

