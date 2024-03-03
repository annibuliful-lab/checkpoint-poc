package stationdevice

import (
	"checkpoint/auth"
	stationDeviceHealthCheck "checkpoint/modules/station-device-health-check"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
)

type StationDeviceResolver struct{}

var stationDeviceService = StationDeviceService{}
var stationDeviceHealthCheckService = stationDeviceHealthCheck.StationDeviceHealthCheckActivityService{}

func (parent StationDevice) HealthActivities(input GetActivitiesInput) ([]*stationDeviceHealthCheck.StationDeviceHealthCheckActivity, error) {
	filter := stationDeviceHealthCheck.GetStationDeviceHealthCheckActivitiesData{
		Limit:           int64(input.Limit),
		Skip:            int64(input.Skip),
		StationDeviceId: uuid.MustParse(string(parent.ID)),
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

func (StationDeviceResolver) GetStationDevices(ctx context.Context, input GetStationDevicesInput) ([]*StationDevice, error) {
	stationDevices, err := stationDeviceService.FindMany(GetStationDevicesData{
		StationLocationId: uuid.MustParse(string(input.StationLocationId)),
		Search:            input.Search,
		Limit:             int64(input.Limit),
		Skip:              int64(input.Skip),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationDevices, nil
}

func (StationDeviceResolver) GetStationDeviceById(ctx context.Context, input GetStationDeviceInput) (*StationDevice, error) {
	stationDevice, err := stationDeviceService.FindById(GetStationDeviceData{
		ID: uuid.MustParse(string(input.ID)),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationDevice, nil
}

func (StationDeviceResolver) UpdateStationDevice(ctx context.Context, input UpdateStationDeviceInput) (*StationDevice, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	stationDevice, err := stationDeviceService.Update(UpdateStationDeviceData{
		ID:              uuid.MustParse(string(input.ID)),
		Title:           input.Title,
		SoftwareVersion: input.SoftwareVersion,
		HardwareVersion: input.HardwareVersion,
		UpdatedBy:       authorization.AccountId,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationDevice, nil
}

func (StationDeviceResolver) DeleteStationDevice(ctx context.Context, input DeleteStationDeviceInput) (*utils.DeleteOperation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	result, err := stationDeviceService.Delete(DeleteStationDeviceData{
		ID:        uuid.MustParse(string(input.ID)),
		DeletedBy: authorization.AccountId,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return result, nil
}

func (StationDeviceResolver) CreateStationDevice(ctx context.Context, input CreateStationDeviceInput) (*StationDevice, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	stationDevice, err := stationDeviceService.Create(CreateStationDeviceData{
		StationLocationId: uuid.MustParse(string(input.StationLocationId)),
		Title:             input.Title,
		HardwareVersion:   input.HardwareVersion,
		SoftwareVersion:   input.SoftwareVersion,
		CreatedBy:         authorization.AccountId,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationDevice, nil
}
