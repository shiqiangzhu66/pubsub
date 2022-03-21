package main

type Publisher interface {
	AddSubscriberCh() chan<- Subscriber
	RemoveSubscriberCh() chan<- Subscriber
	PublishingCh() chan<- interface{}
	Stop()
}

// InMemoryPublisher
type InMemoryPublisher struct {
	subscribers []Subscriber
	addSubCh    chan Subscriber
	removeSubCh chan Subscriber
	in          chan interface{}
	stop        chan struct{}
}

// AddSubscriberCh
//  @receiver p
//  @return chan
func (p *InMemoryPublisher) AddSubscriberCh() chan<- Subscriber {
	return p.addSubCh
}

// RemoveSubscriberCh
//  @receiver p
//  @return chan c
func (p *InMemoryPublisher) RemoveSubscriberCh() chan<- Subscriber {
	return p.removeSubCh
}

// PublishingCh
//  @receiver p
//  @return chan
func (p *InMemoryPublisher) PublishingCh() chan<- interface{} {
	return p.in
}

// Stop
//  @receiver p
func (p *InMemoryPublisher) Stop() {
	close(p.stop)
}

// start
//  @receiver p
func (p *InMemoryPublisher) start() {
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
				if candidate != sub {
					continue
				}

				p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
				candidate.Close()
				break
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

// NewPublisher
//  @return Publisher
func NewPublisher(subs ...Subscriber) Publisher {
	p := &InMemoryPublisher{
		addSubCh:    make(chan Subscriber),
		removeSubCh: make(chan Subscriber),
		in:          make(chan interface{}),
		stop:        make(chan struct{}),
	}

	go p.start()

	for _, sub := range subs {
		p.addSubCh <- sub
	}

	return p
}
