package enum

import (
	"fmt"
)

type BlacklistPriority int

const (
	BLACKLIST_NONE    BlacklistPriority = iota // default value
	BLACKLIST_WARNING BlacklistPriority = 1
	BLACKLIST_DANGER  BlacklistPriority = 2
	BLACKLIST_NORMAL  BlacklistPriority = 3
)

var blacklistPriorityStates = [...]string{"WARNING", "DANGER", "NORMAL", "NONE"}

func GetBlacklistPriority(str string) BlacklistPriority {

	for i, st := range blacklistPriorityStates {
		if st == str {
			return BlacklistPriority(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s BlacklistPriority) String() string { return blacklistPriorityStates[s] }

func (s *BlacklistPriority) Deserialize(str string) {
	var found bool
	for i, st := range blacklistPriorityStates {
		if st == str {
			found = true
			(*s) = BlacklistPriority(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (BlacklistPriority) ImplementsGraphQLType(name string) bool {
	return name == "BlacklistPriority"
}

func (s *BlacklistPriority) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
