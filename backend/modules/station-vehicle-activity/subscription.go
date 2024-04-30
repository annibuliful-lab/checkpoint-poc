package stationvehicleactivity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

func (r *StationVehicleActivityResolver) SetupSubscription() *StationVehicleActivityResolver {

	r.stationVehicleActivityEvent = make(chan *StationVehicleActivity)
	r.stationVehicleActivitySubscriber = make(chan *StationVehicleActivitySubscriber)

	return r
}

func (r *StationVehicleActivityResolver) OnStationVehicleActityEvent(ctx context.Context, input struct {
	StationId graphql.ID
}) <-chan *StationVehicleActivity {
	c := make(chan *StationVehicleActivity)

	// Create a subscriber and manage it in the BroadcastStationVehicleActivity function
	subscriber := &StationVehicleActivitySubscriber{
		stop:  ctx.Done(),
		event: c,
	}

	r.stationVehicleActivitySubscriber <- subscriber

	return c
}

func (r *StationVehicleActivityResolver) BroadcastStationVehicleActivity() {
	subscribers := map[string]*StationVehicleActivitySubscriber{}
	unsubscribe := make(chan string)

	// NOTE: subscribing and sending events are at odds.
	for {
		select {
		case id := <-unsubscribe:
			delete(subscribers, id)
		case s := <-r.stationVehicleActivitySubscriber:
			id := uuid.NewString()
			subscribers[id] = s
		case e := <-r.stationVehicleActivityEvent:

			for id, s := range subscribers {
				go func(id string, s *StationVehicleActivitySubscriber) {
					select {
					case <-s.stop:
						unsubscribe <- id
						return
					default:
					}

					select {
					case <-s.stop:
						unsubscribe <- id
					case s.event <- e:
					case <-time.After(time.Second):
					}
				}(id, s)
			}
		}
	}
}
