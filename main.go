package main

import (
	"time"
)

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
