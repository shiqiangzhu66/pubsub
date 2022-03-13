package main

import (
	"fmt"
	"time"
)

type Subscriber interface {
	Notify(msg interface{}) error
	Close()
}

type Publisher interface {
	start()
	AddSubscriberCh() chan<- Subscriber
	RemoveSubscriberCh() chan<- Subscriber
	PublishingCh() chan<- interface{}
	Stop()
}

// implement
type subscriber struct {
	in chan interface{}
	id int
}

func (s *subscriber) Notify(msg interface{}) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%#v", rec)
		}
	}()

	select {
	case s.in <- msg:
	case <-time.After(time.Second):
		err = fmt.Errorf("timeout")
	}

	return
}

func (s *subscriber) Close() {
	close(s.in)
}

func NewSubscriber(id int) Subscriber {
	s := &subscriber{
		id: id,
		in: make(chan interface{}),
	}

	go func() {
		for msg := range s.in {
			fmt.Printf("(W%d): %v\n", s.id, msg)
		}
	}()

	return s
}

// publisher
type publisher struct {
	subscribers []Subscriber
	addSubCh    chan Subscriber
	removeSubCh chan Subscriber
	in          chan interface{}
	stop        chan struct{}
}

func (p *publisher) AddSubscriberCh() chan<- Subscriber {
	return p.addSubCh
}
func (p *publisher) RemoveSubscriberCh() chan<- Subscriber {
	return p.removeSubCh
}
func (p *publisher) PublishingCh() chan<- interface{} {
	return p.in
}

func (p *publisher) Stop() {
	close(p.stop)
}

func (p *publisher) start() {
	for {
		select {
		case msg := <-p.in:
			for _, sub := range p.subscribers {
				sub.Notify(msg)
			}

		case sub := <-p.addSubCh:
			p.subscribers = append(p.subscribers, sub)

		case sub := <-p.removeSubCh:
			for i, candidate := range p.subscribers {
				if candidate == sub {
					p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
					candidate.Close()
					break
				}
			}
		case <-p.stop:
			for _, sub := range p.subscribers {
				sub.Close()
			}

			close(p.addSubCh)
			close(p.removeSubCh)
			close(p.in)
			return
		}
	}
}

func NewPublisher() Publisher {
	return &publisher{
		addSubCh:    make(chan Subscriber),
		removeSubCh: make(chan Subscriber),
		in:          make(chan interface{}),
		stop:        make(chan struct{}),
	}
}

func main() {
	s1 := NewSubscriber(1)
	s2 := NewSubscriber(2)

	p := NewPublisher()
	go p.start()

	p.AddSubscriberCh() <- s1
	p.AddSubscriberCh() <- s2

	p.PublishingCh() <- "hello"

	p.RemoveSubscriberCh() <- s1
	p.PublishingCh() <- "word"

	time.Sleep(1 * time.Second)
	p.Stop()
}
