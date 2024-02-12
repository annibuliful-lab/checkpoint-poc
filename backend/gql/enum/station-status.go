package enum

import (
	"fmt"
)

type StationStatus int

const (
	STATION_ONLINE StationStatus = iota
	STATION_OFFLINE
)

var StationStatusStates = [...]string{"ONLINE", "OFFLINE"}

func GetStationStatus(str string) StationStatus {

	for i, st := range DeviceStatusStates {
		if st == str {
			return StationStatus(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s StationStatus) String() string { return DeviceStatusStates[s] }

func (s *StationStatus) Deserialize(str string) {
	var found bool
	for i, st := range DeviceStatusStates {
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
