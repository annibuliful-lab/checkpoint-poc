package enum

import (
	"fmt"
)

type DevicePermittedLabel int

const (
	DEVICE_NONE      DevicePermittedLabel = iota // default value
	DEVICE_WHITELIST DevicePermittedLabel = 2
	DEVICE_BLACKLIST DevicePermittedLabel = 1
)

var states = [...]string{"WHITELIST", "BLACKLIST", "NONE"}

func GetDevicePermittedLabel(str string) DevicePermittedLabel {

	for i, st := range states {
		if st == str {
			return DevicePermittedLabel(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s DevicePermittedLabel) String() string { return states[s] }

func (s *DevicePermittedLabel) Deserialize(str string) {
	var found bool
	for i, st := range states {
		if st == str {
			found = true
			(*s) = DevicePermittedLabel(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (DevicePermittedLabel) ImplementsGraphQLType(name string) bool {
	return name == "DevicePermittedLabel"
}

func (s *DevicePermittedLabel) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
