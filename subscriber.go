package main

import (
	"fmt"
	"time"
)

// Subscriber
type Subscriber interface {
	Notify(msg interface{}) error
	Close()
}

// Handler
type Handler interface {
	Handle(msg interface{})
}

// HandlerFunc
//  @param msg
type HandlerFunc func(msg interface{})

// Handle
//  @receiver handler
//  @param msg
func (handler HandlerFunc) Handle(msg interface{}) {
	handler(msg)
}

//
// implement subscriber
//

// InMemorySubscriber
type InMemorySubscriber struct {
	in      chan interface{}
	name    string
	handler Handler
}

// SubscriberOption
//  @param s
type SubscriberOption func(s *InMemorySubscriber)

// WithSubscriberHandler
//  @param handler
//  @return SubscriberOption
func WithSubscriberHandler(handler Handler) SubscriberOption {
	return func(s *InMemorySubscriber) {
		s.handler = handler
	}
}

// WithSubscriberHandlerFunc
//  @param handler
//  @return SubscriberOption
func WithSubscriberHandlerFunc(handler func(msg interface{})) SubscriberOption {
	return func(s *InMemorySubscriber) {
		s.handler = HandlerFunc(handler)
	}
}

// Notify
//  @receiver s
//  @param msg
//  @return err
func (s *InMemorySubscriber) Notify(msg interface{}) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%#v", rec)
		}
	}()

	select {
	case s.in <- msg:
	case <-time.After(1 * time.Second):
		err = fmt.Errorf("timeout")
	}

	return
}

// Close
//  @receiver s
func (s *InMemorySubscriber) Close() {
	close(s.in)
}

// NewSubscriber
//  @param id
//  @return Subscriber
func NewSubscriber(name string, opts ...SubscriberOption) Subscriber {
	s := &InMemorySubscriber{
		name:    name,
		in:      make(chan interface{}, 4),
		handler: HandlerFunc(func(msg interface{}) {}),
	}

	for _, opt := range opts {
		opt(s)
	}

	go func() {
		for msg := range s.in {
			fmt.Printf("(subscriber-%v): %v\n", s.name, msg)
			s.handler.Handle(msg)
		}
	}()

	return s
}
