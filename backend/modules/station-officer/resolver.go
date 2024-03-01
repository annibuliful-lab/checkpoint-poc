package stationofficer

import (
	"checkpoint/auth"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
)

type StationOfficerResolver struct{}

var stationOfficerService = StationOfficerService{}

func (StationOfficerResolver) DeleteStationOfficer(ctx context.Context, input DeleteStationOfficerInput) (*utils.DeleteOperation, error) {
	result, err := stationOfficerService.Delete(DeleteStationOfficerData{
		ID: uuid.MustParse(string(input.ID)),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return result, nil
}

func (StationOfficerResolver) GetStationOfficers(ctx context.Context, input GetStationOfficersInput) ([]*StationOfficer, error) {
	stationOfficers, err := stationOfficerService.FindMany(GetStationOfficersData{
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

	return stationOfficers, nil
}

func (StationOfficerResolver) GetStationOfficerById(ctx context.Context, input GetStationOfficerInput) (*StationOfficer, error) {

	stationOfficer, err := stationOfficerService.FindById(GetStationOfficerData{
		ID: uuid.MustParse(string(input.ID)),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationOfficer, nil
}

func (StationOfficerResolver) UpdateStationOfficer(ctx context.Context, input UpdateStationOfficerInput) (*StationOfficer, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	stationOfficer, err := stationOfficerService.Update(UpdateStationOfficerData{
		ID:        uuid.MustParse(string(input.ID)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
		Firstname: input.Firstname,
		Lastname:  input.Lastname,
		Msisdn:    input.Msisdn,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationOfficer, nil
}

func (StationOfficerResolver) CreateStationOfficer(ctx context.Context, input CreateStationOfficerInput) (*StationOfficer, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	stationOfficer, err := stationOfficerService.Create(CreateStationOfficerData{
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		StationLocationId: uuid.MustParse(string(input.StationLocationId)),
		Firstname:         input.Firstname,
		Lastname:          input.Lastname,
		Msisdn:            input.Msisdn,
		CreatedBy:         authorization.AccountId,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationOfficer, nil
}
