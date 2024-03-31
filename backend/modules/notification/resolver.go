package notification

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type NotificationResolver struct {
	HelloSaidEvents     chan *HelloSaidEvent
	helloSaidSubscriber chan *helloSaidSubscriber
}

func (r *NotificationResolver) SayHello(args struct{ Msg string }) *HelloSaidEvent {
	e := &HelloSaidEvent{Msg: args.Msg, ID: uuid.New().String()}
	go func() {
		select {
		case r.HelloSaidEvents <- e:
		case <-time.After(1 * time.Second):
		}
	}()
	return e
}

func (r *NotificationResolver) WsResolver() *NotificationResolver {
	r.HelloSaidEvents = make(chan *HelloSaidEvent)
	r.helloSaidSubscriber = make(chan *helloSaidSubscriber)

	go r.broadcastHelloSaid()

	return r
}

type helloSaidSubscriber struct {
	stop   <-chan struct{}
	events chan<- *HelloSaidEvent
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
		case e := <-r.HelloSaidEvents:
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

func (r *NotificationResolver) HelloSaid(ctx context.Context) <-chan *HelloSaidEvent {

	c := make(chan *HelloSaidEvent)
	// NOTE: this could take a while
	r.helloSaidSubscriber <- &helloSaidSubscriber{events: c, stop: ctx.Done()}

	return c
}

func (NotificationResolver) OnHello(ctx context.Context) <-chan string {
	c := make(chan string)
	c <- "Hello"
	return c
}
