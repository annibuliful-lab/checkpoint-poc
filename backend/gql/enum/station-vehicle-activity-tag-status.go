package enum

import (
	"fmt"
)

type StationVehicleActivityTagStatus int

const (
	VEHICLE_ACTIVITY_LICENSE_PLATE StationVehicleActivityTagStatus = iota
	VEHICLE_ACTIVITY_IMEI
	VEHICLE_ACTIVITY_IMSI
)

var StationVehicleActivityTagStatusStates = [...]string{
	"LICENSE_PLATE",
	"IMEI",
	"IMSI",
}

func GetStationVehicleActivityTagStatus(str string) StationVehicleActivityTagStatus {

	for i, st := range StationVehicleActivityTagStatusStates {
		if st == str {
			return StationVehicleActivityTagStatus(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s StationVehicleActivityTagStatus) String() string {
	return StationVehicleActivityTagStatusStates[s]
}

func (s *StationVehicleActivityTagStatus) Deserialize(str string) {
	var found bool
	for i, st := range StationVehicleActivityTagStatusStates {
		if st == str {
			found = true
			(*s) = StationVehicleActivityTagStatus(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (StationVehicleActivityTagStatus) ImplementsGraphQLType(name string) bool {
	return name == "StationVehicleActivityTagStatus"
}

func (s *StationVehicleActivityTagStatus) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
