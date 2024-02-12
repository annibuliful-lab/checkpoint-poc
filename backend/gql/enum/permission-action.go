package enum

import (
	"fmt"
)

type PermissionAction int

const (
	CREATE PermissionAction = iota
	UPDATE
	DELETE
	READ
)

var permissionActionStates = [...]string{"CREATE", "UPDATE", "DELETE", "READ"}

func GetPermissionAction(str string) PermissionAction {

	for i, st := range permissionActionStates {
		if st == str {
			return PermissionAction(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s PermissionAction) String() string { return permissionActionStates[s] }

func (s *PermissionAction) Deserialize(str string) {
	var found bool
	for i, st := range DeviceStatusStates {
		if st == str {
			found = true
			(*s) = PermissionAction(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (PermissionAction) ImplementsGraphQLType(name string) bool {
	return name == "PermissionAction"
}

func (s *PermissionAction) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
