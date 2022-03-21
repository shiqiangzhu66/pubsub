package main

import (
	"time"
)

func main() {
	s1 := NewSubscriber("s1")
	s2 := NewSubscriber("s2")
	s3 := NewSubscriber("s3")

	p := NewPublisher(s1, s2, s3)

	// p.AddSubscriberCh() <- s1
	// p.AddSubscriberCh() <- s2
	// p.AddSubscriberCh() <- s3

	p.PublishingCh() <- "hello"

	p.RemoveSubscriberCh() <- s2
	p.PublishingCh() <- "word"

	time.Sleep(1 * time.Second)
	p.Stop()
}
