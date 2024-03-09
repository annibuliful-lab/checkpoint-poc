package notification

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type NotificationResolver struct {
	helloSaidEvents     chan *helloSaidEvent
	helloSaidSubscriber chan *helloSaidSubscriber
}

func (r *NotificationResolver) SayHello(args struct{ Msg string }) *helloSaidEvent {
	e := &helloSaidEvent{Msg: args.Msg, ID: uuid.New().String()}
	go func() {
		select {
		case r.helloSaidEvents <- e:
		case <-time.After(1 * time.Second):
		}
	}()
	return e
}

func (r *NotificationResolver) WsResolver() *NotificationResolver {
	r.helloSaidEvents = make(chan *helloSaidEvent)
	r.helloSaidSubscriber = make(chan *helloSaidSubscriber)

	go r.broadcastHelloSaid()

	return r
}

type helloSaidSubscriber struct {
	stop   <-chan struct{}
	events chan<- *helloSaidEvent
}

func (r *NotificationResolver) broadcastHelloSaid() {
	subscribers := map[string]*helloSaidSubscriber{}
	unsubscribe := make(chan string)

	// NOTE: subscribing and sending events are at odds.
	for {
		select {
		case id := <-unsubscribe:
			delete(subscribers, id)
		case s := <-r.helloSaidSubscriber:
			id := uuid.New().String()
			subscribers[id] = s
		case e := <-r.helloSaidEvents:
			for id, s := range subscribers {
				go func(id string, s *helloSaidSubscriber) {
					select {
					case <-s.stop:
						unsubscribe <- id
						return
					default:
					}

					select {
					case <-s.stop:
						unsubscribe <- id
					case s.events <- e:
					case <-time.After(time.Second):
					}
				}(id, s)
			}
		}
	}
}

func (r *NotificationResolver) HelloSaid(ctx context.Context) <-chan *helloSaidEvent {

	c := make(chan *helloSaidEvent)
	// NOTE: this could take a while
	r.helloSaidSubscriber <- &helloSaidSubscriber{events: c, stop: ctx.Done()}

	return c
}

type helloSaidEvent struct {
	ID  string
	Msg string
}
