package enum

import (
	"fmt"
)

type RemarkState int

const (
	REMARK_NONE    RemarkState = iota // default value
	REMARK_WARNING RemarkState = 1
	REMARK_DANGER  RemarkState = 2
	REMARK_NORMAL  RemarkState = 3
)

var remarkStates = [...]string{
	"WHITELIST",
	"BLACKLIST",
	"IN_QUEUE",
	"IN_PROGRESS",
	"PASSED",
	"WAITING",
	"INVESTIGATING",
	"SUSPICION",
}

func GetRemarkState(str string) RemarkState {

	for i, st := range remarkStates {
		if st == str {
			return RemarkState(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s RemarkState) String() string { return remarkStates[s] }

func (s *RemarkState) Deserialize(str string) {
	var found bool
	for i, st := range remarkStates {
		if st == str {
			found = true
			(*s) = RemarkState(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (RemarkState) ImplementsGraphQLType(name string) bool {
	return name == "RemarkState"
}

func (s *RemarkState) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
