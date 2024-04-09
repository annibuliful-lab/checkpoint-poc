package enum

import (
	"fmt"
)

type ImageType int

const (
	FRONT ImageType = iota
	REAR
	DRIVER
	LICENSE_PLATE
	NONE
	CONFIG
)

var ImageTypeStates = [...]string{"FRONT", "REAR", "DRIVER", "LICENSE_PLATE", "NONE", "CONFIG"}

func GetImageType(str string) ImageType {

	for i, st := range ImageTypeStates {
		if st == str {
			return ImageType(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s ImageType) String() string { return ImageTypeStates[s] }

func (s *ImageType) Deserialize(str string) {
	var found bool
	for i, st := range ImageTypeStates {
		if st == str {
			found = true
			(*s) = ImageType(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (ImageType) ImplementsGraphQLType(name string) bool {
	return name == "ImageType"
}

func (s *ImageType) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
