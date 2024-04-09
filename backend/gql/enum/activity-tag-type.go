package enum

import (
	"fmt"
)

type ActivityTagType int

const (
	ACTIVITY_LICENSE_PLATE ActivityTagType = iota
	ACTIVITY_IMEI
	ACTIVITY_IMSI
)

var ActivityTagTypeState = [...]string{
	"LICENSE_PLATE",
	"IMEI",
	"IMSI",
}

func GetActivityTagType(str string) ActivityTagType {

	for i, st := range ActivityTagTypeState {
		if st == str {
			return ActivityTagType(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s ActivityTagType) String() string {
	return ActivityTagTypeState[s]
}

func (s *ActivityTagType) Deserialize(str string) {
	var found bool
	for i, st := range ActivityTagTypeState {
		if st == str {
			found = true
			(*s) = ActivityTagType(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (ActivityTagType) ImplementsGraphQLType(name string) bool {
	return name == "ActivityTagType"
}

func (s *ActivityTagType) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
