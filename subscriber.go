package main

import (
	"fmt"
	"time"
)

type Subscriber interface {
	Notify(msg interface{}) error
	Close()
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
