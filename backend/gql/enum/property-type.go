package enum

import (
	"fmt"
)

type PropertyType int

var propertyTypeStates = [...]string{
	"VEHICLE_COLOR",
	"LP_TYPE",
	"VEHICLE_BRAND",
	"VEHICLE_TYPE",
	"VEHICLE_MODEL",
}

func GetPropertyType(str string) PropertyType {

	for i, st := range permissionActionStates {
		if st == str {
			return PropertyType(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s PropertyType) String() string { return propertyTypeStates[s] }

func (s *PropertyType) Deserialize(str string) {
	var found bool
	for i, st := range propertyTypeStates {
		if st == str {
			found = true
			(*s) = PropertyType(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (PropertyType) ImplementsGraphQLType(name string) bool {
	return name == "PropertyType"
}

func (s *PropertyType) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
