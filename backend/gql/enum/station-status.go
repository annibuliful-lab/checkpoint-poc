package enum

import (
	"fmt"
)

type StationStatus int

const (
	STATION_ONLINE StationStatus = iota
	STATION_OFFLINE
	STATION_CLOSED
	STATION_MAINTENANCE
)

var StationStatusStates = [...]string{"ONLINE", "OFFLINE", "CLOSED", "MAINTENANCE"}

func GetStationStatus(str string) StationStatus {

	for i, st := range StationStatusStates {
		if st == str {
			return StationStatus(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s StationStatus) String() string { return StationStatusStates[s] }

func (s *StationStatus) Deserialize(str string) {
	var found bool
	for i, st := range StationStatusStates {
		if st == str {
			found = true
			(*s) = StationStatus(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (StationStatus) ImplementsGraphQLType(name string) bool {
	return name == "StationStatus"
}

func (s *StationStatus) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
