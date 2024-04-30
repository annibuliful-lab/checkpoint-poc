package stationvehicleactivity

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

func (r StationVehicleActivityResolver) SetupSubscription() StationVehicleActivityResolver {

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

	// Ensure we unsubscribe when the context is cancelled
	go func() {
		<-ctx.Done()
		close(c)
	}()

	return c
}

func (r *StationVehicleActivityResolver) BroadcastStationVehicleActivity() {
	subscribers := make(map[*StationVehicleActivitySubscriber]struct{})

	// Handle incoming events and manage subscribers
	for {
		select {
		case s := <-r.stationVehicleActivitySubscriber:
			subscribers[s] = struct{}{}
		case e := <-r.stationVehicleActivityEvent:
			for s := range subscribers {
				select {
				case <-s.stop:
					delete(subscribers, s)
				case s.event <- e:
				}
			}
		}
	}
}
