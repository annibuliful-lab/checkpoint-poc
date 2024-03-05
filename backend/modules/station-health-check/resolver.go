package stationhealthcheck

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/auth"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
)

type StationLocationHealthCheckActivityResolver struct{}

var stationLocationHealthCheckActivityService = StationHealthCheckActivityService{}

func (StationLocationHealthCheckActivityResolver) UpdateStationLocationHealthCheckActivity(ctx context.Context, input UpdateStationHealthCheckActivityInput) (*StationLocationHealthCheckActivity, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	fieldsToUpdate := UpdateStationHealthCheckActivityData{
		UpdatedBy: authorization.AccountId,
		ID:        uuid.MustParse(string(input.ID)),
	}

	if input.StationStatus != nil {
		stationStatus := model.StationStatus(input.StationStatus.String())
		fieldsToUpdate.StationStatus = &stationStatus
	}

	if input.EndDatetime != nil {
		fieldsToUpdate.EndDatetime = &input.EndDatetime.Time
	}

	if input.StartDatetime != nil {
		fieldsToUpdate.StartDatetime = &input.StartDatetime.Time
	}

	activity, err := stationLocationHealthCheckActivityService.Update(fieldsToUpdate)

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return activity, nil
}

func (StationLocationHealthCheckActivityResolver) GetStationLocationHealthCheckActivities(ctx context.Context, input GetStationHealthCheckActivitiesInput) ([]*StationLocationHealthCheckActivity, error) {
	filter := GetStationHealthCheckActivitiesData{
		StationId: uuid.MustParse(string(input.StationId)),
		Limit:     int64(input.Limit),
		Skip:      int64(input.Skip),
	}

	if input.EndDatetime != nil {
		filter.EndDatetime = &input.EndDatetime.Time
	}

	if input.StartDatetime != nil {
		filter.StartDatetime = &input.StartDatetime.Time
	}

	if input.StationStatus != nil {
		status := model.StationStatus(input.StationStatus.String())
		filter.StationStatus = &status
	}

	activities, err := stationLocationHealthCheckActivityService.FindMany(filter)

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return activities, nil
}

func (StationLocationHealthCheckActivityResolver) CreateStationLocationHealthCheckActivity(ctx context.Context, input CreateStationHealthCheckActivityInput) (*StationLocationHealthCheckActivity, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	fieldsToCreate := CreateStationHealthCheckActivityData{
		CreatedBy:     authorization.AccountId,
		StationId:     uuid.MustParse(string(input.StationId)),
		StationStatus: model.StationStatus(input.StationStatus.String()),
		StartDatetime: input.StartDatetime.Time,
	}

	if input.EndDatetime != nil {
		fieldsToCreate.EndDatetime = &input.EndDatetime.Value.Time
	}

	activity, err := stationLocationHealthCheckActivityService.Create(fieldsToCreate)

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return activity, nil
}
