package vehicletargetconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/auth"
	"checkpoint/modules/tag"
	"checkpoint/modules/upload"
	vehicletargetconfigurationimage "checkpoint/modules/vehicle-target-configuration-image"
	"checkpoint/utils"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

var vehicleService = VehicleTargetConfigurationService{}
var tagDataloader = tag.TagService{}.VehicleTargetConfigurationTags()
var vehicleTargetImageService = vehicletargetconfigurationimage.VehicleTargetConfigurationImageService{}
var imageDatalaoder = vehicleTargetImageService.Dataloader()

type VehicleTargetConfigurationResolver struct{}

func (VehicleTargetConfigurationResolver) CreateVehicleTargetConfiguration(ctx context.Context, input CreateVehicleTargetConfigurationInput) (*VehicleTargetConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	vehicleTarget, err := vehicleService.Create(CreateVehicleTargetConfigurationData{
		StationLocationId: uuid.MustParse(string(input.StationLocationId)),
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		CreatedBy:         authorization.AccountId,
		Prefix:            input.Prefix,
		Number:            input.Number,
		Province:          input.Province,
		Type:              input.Type,
		Brand:             input.Brand,
		Color:             input.Color,
		Country:           input.Country,
		PermittedLabel:    model.DevicePermittedLabel(input.PermittedLabel.String()),
		BlacklistPriority: model.BlacklistPriority(input.BlacklistPriority.String()),
		Tags:              input.Tags,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return vehicleTarget, nil
}

func (VehicleTargetConfigurationResolver) UpdateVehicleTargetConfiguration(ctx context.Context, input UpdateVehicleTargetConfigurationInput) (*VehicleTargetConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	fieldsToUpdate := UpdateVehicleTargetConfigurationData{
		ID:        uuid.MustParse(string(input.ID)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
		UpdatedBy: authorization.AccountId,
		Prefix:    input.Prefix,
		Number:    input.Number,
		Province:  input.Province,
		Type:      input.Type,
		Brand:     input.Brand,
		Color:     input.Color,
		Country:   input.Country,
		Tags:      input.Tags,
	}

	if input.BlacklistPriority != nil {
		blacklist := model.BlacklistPriority(input.BlacklistPriority.String())
		fieldsToUpdate.BlacklistPriority = &blacklist
	}

	if input.PermittedLabel != nil {
		label := model.DevicePermittedLabel(input.PermittedLabel.String())
		fieldsToUpdate.PermittedLabel = &label
	}

	vehicleTarget, err := vehicleService.Update(fieldsToUpdate)

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return vehicleTarget, nil
}

func (VehicleTargetConfigurationResolver) DeleteVehicleTargetConfiguration(ctx context.Context, input DeleteVehicleTargetConfigurationInput) (*utils.DeleteOperation, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	result, err := vehicleService.Delete(DeleteVehicleTargetConfigurationData{
		ID:        uuid.MustParse(string(input.ID)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
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

func (VehicleTargetConfigurationResolver) GetVehicleTargetConfigurationById(ctx context.Context, input GetVehicleTargetConfigurationInput) (*VehicleTargetConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	vehicleTarget, err := vehicleService.FindById(GetVehicleTargetConfigurationData{
		ID:        uuid.MustParse(string(input.ID)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return vehicleTarget, nil
}

func (VehicleTargetConfigurationResolver) GetVehicleTargetConfigurations(ctx context.Context, input GetVehicleTargetsConfigurationInput) ([]*VehicleTargetConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	filter := GetVehicleTargetsConfigurationData{
		ProjectId: uuid.MustParse(authorization.ProjectId),
		Type:      input.Type,
		Search:    input.Search,
		Tags:      input.Tags,
		Limit:     int64(input.Limit),
		Skip:      int64(input.Skip),
	}

	if input.StationLocationId != nil {
		stationLocationId := uuid.MustParse(string(*input.StationLocationId))
		filter.StationLocationId = &stationLocationId
	}

	if input.BlacklistPriority != nil {
		blacklist := model.BlacklistPriority(input.BlacklistPriority.String())
		filter.BlacklistPriority = &blacklist
	}

	if input.PermittedLabel != nil {
		label := model.DevicePermittedLabel(input.PermittedLabel.String())
		filter.PermittedLabel = &label
	}

	vehicleTargets, err := vehicleService.FindMany(filter)

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return vehicleTargets, nil
}

func (parent VehicleTargetConfiguration) Tags(ctx context.Context) (*[]tag.Tag, error) {
	thunk := tagDataloader.Load(ctx, dataloader.StringKey(parent.ID))
	tagsLoaderResult, err := thunk()
	if err != nil {
		log.Println(err.Error())
		return nil, nil
	}

	tags := lo.Map(tagsLoaderResult.([]tag.VehicleTargetConfigurationTag), func(item tag.VehicleTargetConfigurationTag, index int) tag.Tag {
		return tag.Tag{
			Id:        graphql.ID(item.Tag.ID.String()),
			Title:     item.Tag.Title,
			ProjectId: graphql.ID(item.Tag.ProjectId.String()),
		}
	})

	return &tags, nil
}

func (parent VehicleTargetConfiguration) Images(ctx context.Context) (*[]vehicletargetconfigurationimage.VehicleTargetConfigurationImage, error) {
	thunk := imageDatalaoder.Load(ctx, dataloader.StringKey(parent.ID))
	tagsLoaderResult, err := thunk()
	if err != nil {
		log.Println(err.Error())
		return nil, nil
	}

	images := lo.Map(tagsLoaderResult.([]model.VehicleTargetConfigurationImage), func(item model.VehicleTargetConfigurationImage, index int) vehicletargetconfigurationimage.VehicleTargetConfigurationImage {
		url, err := upload.SignedUrl(ctx, item.S3key)
		if err != nil {
			log.Println("get-signed-url-error", err.Error())
		}

		return vehicletargetconfigurationimage.VehicleTargetConfigurationImage{
			Id:                           graphql.ID(item.ID.String()),
			VehicleTargetConfigurationId: graphql.ID(item.VehicleTargetConfigurationId.String()),
			Url:                          *url,
			S3Key:                        item.S3key,
		}
	})

	return &images, nil
}
