package stationdevicehealthcheck

import (
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
)

type StationDeviceHealthCheckActivityResolver struct{}

var stationDeviceHealthCheckService = StationDeviceHealthCheckActivityService{}

func (StationDeviceHealthCheckActivityResolver) CreateStationDeviceHealthCheckActivity(ctx context.Context, input CreateStationDeviceHealthCheckActivityInput) (*StationDeviceHealthCheckActivity, error) {
	activity, err := stationDeviceHealthCheckService.Create(CreateStationDeviceHealthCheckActivityData{
		Issue:           input.Issue,
		StationDeviceId: uuid.MustParse(string(input.StationDeviceId)),
		ActivityTime:    input.ActivityTime.Time,
		Status:          input.Status,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return activity, nil
}

func (StationDeviceHealthCheckActivityResolver) GetStationDeviceHealthCheckActivities(ctx context.Context, input GetStationDeviceHealthCheckActivitiesInput) ([]*StationDeviceHealthCheckActivity, error) {

	filter := GetStationDeviceHealthCheckActivitiesData{
		Limit:           int64(input.Limit),
		Skip:            int64(input.Skip),
		StationDeviceId: uuid.MustParse(string(input.StationDeviceId)),
		Status:          input.Status,
	}

	if input.EndDate.Value != nil {
		filter.EndDate = &input.EndDate.Value.Time
	}

	if input.StartDate.Value != nil {
		filter.StartDate = &input.StartDate.Value.Time
	}

	activities, err := stationDeviceHealthCheckService.FindMany(filter)

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return activities, nil
}
