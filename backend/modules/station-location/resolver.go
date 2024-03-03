package stationlocation

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/auth"
	stationdevice "checkpoint/modules/station-device"
	stationofficer "checkpoint/modules/station-officer"
	"checkpoint/modules/tag"
	"checkpoint/utils"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type StationLocationResolver struct{}

var stationLocationService = StationLocationService{}
var tagService = tag.TagService{}
var stationOfficerService = stationofficer.StationOfficerService{}
var stationDeviceService = stationdevice.StationDeviceService{}
var tagDataloader = tagService.StationLocationTagDataloader()
var stationOfficerLoader = stationOfficerService.StationLocationOfficerDataloader()
var stationDeviceLoader = stationDeviceService.StationLocationDataloader()

func (StationLocationResolver) DeleteStationLocation(ctx context.Context, input DeleteStationLocationInput) (*utils.DeleteOperation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(authorization.ProjectId),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	err := stationLocationService.Delete(DeleteStationLocationData{
		ID:        uuid.MustParse(string(input.ID)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return &utils.DeleteOperation{
		Success: true,
	}, nil
}

func (StationLocationResolver) GetStationLocations(ctx context.Context, input GetStationLocationsInput) ([]*StationLocation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	stationLocations, err := stationLocationService.FindMany(GetStationLocationsData{
		Limit:     int64(input.Limit),
		Skip:      int64(input.Skip),
		Search:    input.Search,
		ProjectId: uuid.MustParse(authorization.ProjectId),
		Tags:      input.Tags,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationLocations, nil
}

func (StationLocationResolver) GetStationLocationById(ctx context.Context, input GetStationLocationByIdInput) (*StationLocation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	match := auth.VerifyProjectAccount(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(authorization.ProjectId),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	stationLocation, err := stationLocationService.FindById(GetStationLocationByIdData{
		Id:        uuid.MustParse(string(input.Id)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationLocation, nil
}

func (StationLocationResolver) UpdateStationLocation(ctx context.Context, input UpdateStationLocationInput) (*StationLocation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(authorization.ProjectId),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	stationLocation, err := stationLocationService.Update(UpdateStationLocationData{
		UpdatedBy:   authorization.AccountId,
		Id:          uuid.MustParse(string(input.Id)),
		ProjectId:   uuid.MustParse(authorization.ProjectId),
		Title:       input.Title,
		Description: input.Description,
		Department:  input.Department,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
		Remark:      input.Remark,
		Tags:        input.Tags,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationLocation, nil

}

func (StationLocationResolver) CreateStationLocation(ctx context.Context, input CreateStationLocationInput) (*StationLocation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(authorization.ProjectId),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	stationLocation, err := stationLocationService.Create(CreateStationLocationData{
		ProjectId:   authorization.ProjectId,
		Department:  input.Department,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
		Title:       input.Title,
		Description: input.Description,
		CreatedBy:   authorization.AccountId,
		Tags:        input.Tags,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return stationLocation, nil
}

func (parent StationLocation) Tags(ctx context.Context) (*[]tag.Tag, error) {
	thunk := tagDataloader.Load(ctx, dataloader.StringKey(parent.ID))
	tagsLoaderResult, err := thunk()
	if err != nil {
		log.Println(err.Error())
		return nil, nil
	}

	tags := lo.Map(tagsLoaderResult.([]tag.StationLocationTag), func(item tag.StationLocationTag, index int) tag.Tag {
		return tag.Tag{
			Id:        graphql.ID(item.Tag.ID.String()),
			Title:     item.Tag.Title,
			ProjectId: graphql.ID(item.Tag.ProjectId.String()),
		}
	})

	return &tags, nil
}

func (parent StationLocation) Officers(ctx context.Context) (*[]stationofficer.StationOfficer, error) {
	thunk := stationOfficerLoader.Load(ctx, dataloader.StringKey(parent.ID))
	officersLoaderResult, err := thunk()
	if err != nil {
		log.Println(err.Error())
		return nil, nil
	}

	officers := lo.Map(officersLoaderResult.([]model.StationOfficer), func(item model.StationOfficer, index int) stationofficer.StationOfficer {
		return stationofficer.StationOfficer{
			ID:                graphql.ID(item.ID.String()),
			StationLocationId: graphql.ID(item.StationLocationId.String()),
			Firstname:         item.Firstname,
			Lastname:          item.Lastname,
			Msisdn:            item.Msisdn,
		}
	})

	return &officers, nil
}

func (parent StationLocation) Devices(ctx context.Context) (*[]stationdevice.StationDevice, error) {
	thunk := stationDeviceLoader.Load(ctx, dataloader.StringKey(parent.ID))
	devicesLoaderResult, err := thunk()
	if err != nil {
		log.Println(err.Error())
		return nil, nil
	}

	devices := lo.Map(devicesLoaderResult.([]model.StationDevice), func(item model.StationDevice, index int) stationdevice.StationDevice {
		return stationdevice.StationDevice{
			ID:                graphql.ID(item.ID.String()),
			StationLocationId: graphql.ID(item.StationLocationId.String()),
			Title:             item.Title,
			HardwareVersion:   item.HardwareVersion,
			SoftwareVersion:   item.SoftwareVersion,
		}
	})

	return &devices, nil
}
