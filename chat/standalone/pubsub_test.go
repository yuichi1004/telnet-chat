package standalone

import (
	"fmt"
	"sync"
	"testing"
)

func TestPubSub_Subscribe(t *testing.T) {
	type fields struct {
		subs   map[int]chan string
		m      sync.Mutex
		serial int
	}
	type expects struct {
		subs int
		serial int
	}
	type args struct {
		c chan string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		test    func(p *PubSub) error
		expects expects
		wantErr bool
	}{
		{
			name: "normal case",
			fields: fields{
				subs: make(map[int]chan string),
			},
			test: func(p *PubSub) error {
				ch := make(chan string)
				_, err := p.Subscribe(ch)
				return err
			},
			expects: expects {
				subs: 1,
				serial: 1,
			},
		},
		{
			name: "normal case - closed",
			fields: fields{
				subs: make(map[int]chan string),
			},
			test: func(p *PubSub) error {
				ch := make(chan string)
				closer, err := p.Subscribe(ch)
				ch2 := make(chan string)
				_, err = p.Subscribe(ch2)
				closer()
				return err
			},
			expects: expects {
				subs: 1,
				serial: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PubSub{
				subs:   tt.fields.subs,
				m:      tt.fields.m,
				serial: tt.fields.serial,
			}
			err := tt.test(p)
			if (err != nil) != tt.wantErr {
				t.Errorf("PubSub.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(p.subs) != tt.expects.subs {
				t.Errorf("len(PubSub.subs) expected = %d, got %d", tt.expects.subs, len(p.subs))
			}
		})
	}
}

func TestPubSub_Publish(t *testing.T) {
	subscribers := make(map[int]chan string)
	subscribers[0] = make(chan string, 1)
	subscribers[1] = make(chan string, 1)

	type fields struct {
		subs   map[int]chan string
		m      sync.Mutex
		serial int
	}
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		check   func() error
		wantErr bool
	}{
		{
			name: "normal case",
			fields: fields {
				subs: subscribers,
				serial: 2,
			},
			args: args{"hello"},
			check: func() error {
				got := 0
				for i := 0; i < 2; i++ {
					select {
					case m1 := <- subscribers[0]:
						if m1 != "hello" {
							return fmt.Errorf("unexpected message (expects: %s, got: %s)", "hello", m1)
						}
						got = got + 1
					case m2 := <- subscribers[1]:
						if m2 != "hello" {
							return fmt.Errorf("unexpected message (expects: %s, got: %s)", "hello", m2)
						}
						got = got + 1
					default:
					}
				}
				if got != 2 {
					return fmt.Errorf("lack of message recieved (expectes: %d, got: %d)", 2, got)
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PubSub{
				subs:   tt.fields.subs,
				m:      tt.fields.m,
				serial: tt.fields.serial,
			}
			if err := p.Publish(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("PubSub.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.check(); err != nil {
				t.Errorf("%s", err)
			}
		})
	}
}
