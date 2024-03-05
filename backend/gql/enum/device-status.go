package enum

import (
	"fmt"
)

type DeviceStatus int

const (
	DEVICE_ONLINE DeviceStatus = iota
	DEVICE_OFFLINE
)

var DeviceStatusStates = [...]string{"ONLINE", "OFFLINE"}

func GetDeviceStatus(str string) DeviceStatus {

	for i, st := range DeviceStatusStates {
		if st == str {
			return DeviceStatus(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s DeviceStatus) String() string { return DeviceStatusStates[s] }

func (s *DeviceStatus) Deserialize(str string) {
	var found bool
	for i, st := range DeviceStatusStates {
		if st == str {
			found = true
			(*s) = DeviceStatus(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (DeviceStatus) ImplementsGraphQLType(name string) bool {
	return name == "DeviceStatus"
}

func (s *DeviceStatus) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
