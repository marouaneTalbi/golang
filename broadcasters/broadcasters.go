package broadcasters

import "fmt"

type broadcaster struct {
	input   chan interface{}
	reg     chan chan<- interface{}
	unreg   chan chan<- interface{}
	outputs map[chan<- interface{}]bool
}

type Broadcaster interface {
	Register(chan<- interface{})
	Unregister(chan<- interface{})
	Close() error
	Submit(interface{})
	TrySubmit(interface{}) bool
}

func (b *broadcaster) broadcast(m interface{}) {
	for ch := range b.outputs {
		ch <- m
	}
}

func (b *broadcaster) run() {
	for {
		select {
		case m := <-b.input:
			b.broadcast(m)

		case ch, ok := <-b.reg:
			if ok {
				b.outputs[ch] = true
			} else {
				return
			}
		case ch := <-b.unreg:
			delete(b.outputs, ch)
		}
	}
}

func NewBroadcaster(buflen int) Broadcaster {
	b := &broadcaster{
		input:   make(chan interface{}, buflen),
		reg:     make(chan chan<- interface{}),
		unreg:   make(chan chan<- interface{}),
		outputs: make(map[chan<- interface{}]bool),
	}

	go b.run()

	return b
}

func (b *broadcaster) Register(newch chan<- interface{}) {
	b.reg <- newch
}

func (b *broadcaster) Unregister(newch chan<- interface{}) {
	b.unreg <- newch
}

func (b *broadcaster) Close() error {
	close(b.reg)
	close(b.unreg)
	return nil
}

func (b *broadcaster) Submit(m interface{}) {

	fmt.Println("submit")
	if b != nil {
		b.input <- m
	}
}

func (b *broadcaster) TrySubmit(m interface{}) bool {
	if b == nil {
		return false
	}
	select {
	case b.input <- m:
		return true
	default:
		return false
	}
}
