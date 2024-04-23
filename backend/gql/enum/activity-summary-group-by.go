package enum

import (
	"fmt"
)

type VehicleActivitySummaryCategory int

const (
	VEHICLE_ACTIVITY_DAY VehicleActivitySummaryCategory = iota
	VEHICLE_ACTIVITY_WEEK
	VEHICLE_ACTIVITY_MONTH
	VEHICLE_ACTIVITY_YEAR
	VEHICLE_ACTIVITY_CUSTOM_DATE
)

var VehicleActivitySummaryCategoryState = [...]string{
	"DAY",
	"WEEK",
	"MONTH",
	"YEAR",
	"CUSTOM_DATE",
}

func GetVehicleActivitySummaryCategory(str string) VehicleActivitySummaryCategory {

	for i, st := range VehicleActivitySummaryCategoryState {
		if st == str {
			return VehicleActivitySummaryCategory(i)
		}
	}

	panic("invalid value for enum State: " + str)

}

func (s VehicleActivitySummaryCategory) String() string {
	return VehicleActivitySummaryCategoryState[s]
}

func (s *VehicleActivitySummaryCategory) Deserialize(str string) {
	var found bool
	for i, st := range VehicleActivitySummaryCategoryState {
		if st == str {
			found = true
			(*s) = VehicleActivitySummaryCategory(i)
		}
	}
	if !found {
		panic("invalid value for enum State: " + str)
	}
}

func (VehicleActivitySummaryCategory) ImplementsGraphQLType(name string) bool {
	return name == "VehicleActivitySummaryCategory"
}

func (s *VehicleActivitySummaryCategory) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		s.Deserialize(input)
	default:
		err = fmt.Errorf("wrong type for State: %T", input)
	}
	return err
}
