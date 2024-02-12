package mobiledeviceconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/auth"
	imeiconfiguration "checkpoint/modules/imei-configuration"
	imsiconfiguration "checkpoint/modules/imsi-configuration"
	"checkpoint/modules/tag"
	"checkpoint/utils"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type MobileDeviceConfigurationResolver struct{}

var mobileDeviceService = MobileDeviceConfigurationService{}
var imeiService = imeiconfiguration.ImeiConfigurationService{}
var imsiService = imsiconfiguration.ImsiConfigurationService{}
var imsiDataloader = imsiService.Dataloader()
var imeiDataloader = imeiService.Dataloader()
var tagService = tag.TagService{}
var tagDataloader = tagService.MobileDeviceConfigurationTagDataloader()

func (MobileDeviceConfigurationResolver) GetMobileDeviceConfigurations(ctx context.Context, args GetMobileDeviceConfigurationsInput) ([]MobileDeviceConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	mobileDeviceConfigurations, err := mobileDeviceService.FindMany(GetMobileDeviceConfigurationsData{
		StationLocationId: uuid.MustParse(string(args.StationLocationId)),
		Search:            args.Search,
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		Tags:              args.Tags,
		PermittedLabel:    args.PermittedLabel,
		BlacklistPriority: args.BlacklistPriority,
		pagination: utils.OffsetPagination{
			Limit: int64(args.Limit),
			Skip:  int64(args.Skip),
		},
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return mobileDeviceConfigurations, nil
}

func (MobileDeviceConfigurationResolver) GetMobileDeviceConfigurationById(ctx context.Context, args GetMobileDeviceConfigurationInput) (*MobileDeviceConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	mobileDevice, err := mobileDeviceService.FindById(GetMobileDeviceConfigurationData{
		ID:        uuid.MustParse(string(args.ID)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return mobileDevice, nil
}

func (MobileDeviceConfigurationResolver) UpdateMobileDeviceConfiguration(ctx context.Context, args UpdateMobileDeviceConfigurationInput) (*MobileDeviceConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	if args.ReferenceImeiConfigurationId != nil {
		_, _, err := imeiService.FindById(imeiconfiguration.GetImeiConfigurationData{
			ID:        uuid.MustParse(string(*args.ReferenceImeiConfigurationId)),
			ProjectId: uuid.MustParse(authorization.ProjectId),
		})
		if err != nil {
			return nil, utils.GraphqlError{
				Code:    err.Error(),
				Message: err.Error(),
			}
		}
	}

	if args.ReferenceImsiConfigurationId != nil {
		_, _, err := imsiService.FindById(imsiconfiguration.GetImsiConfigurationByIdData{
			ID:        uuid.MustParse(string(*args.ReferenceImsiConfigurationId)),
			ProjectId: uuid.MustParse(authorization.ProjectId),
		})

		if err != nil {
			return nil, utils.GraphqlError{
				Code:    err.Error(),
				Message: err.Error(),
			}
		}
	}

	var mobileDeviceInput = UpdateMobileDeviceConfigurationData{
		ID:                uuid.MustParse(string(args.ID)),
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		PermittedLabel:    args.PermittedLabel,
		BlacklistPriority: args.BlacklistPriority,
		Msisdn:            args.Msisdn,
		Tags:              args.Tags,
		UpdatedBy:         authorization.AccountId,
	}

	if args.ReferenceImeiConfigurationId != nil {
		id := uuid.MustParse(string(*args.ReferenceImeiConfigurationId))
		mobileDeviceInput.ReferenceImeiConfigurationId = &id
	}

	if args.ReferenceImsiConfigurationId != nil {
		id := uuid.MustParse(string(*args.ReferenceImsiConfigurationId))
		mobileDeviceInput.ReferenceImsiConfigurationId = &id
	}

	mobileDevice, err := mobileDeviceService.Update(mobileDeviceInput)

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return mobileDevice, nil
}

func (MobileDeviceConfigurationResolver) DeleteMobileDeviceConfiguration(ctx context.Context, args DeleteMobileDeviceConfigurationInpt) (*utils.DeleteOperation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	result, err := mobileDeviceService.Delete(DeleteMobileDeviceConfigurationData{
		ID:        uuid.MustParse(string(args.ID)),
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

func (MobileDeviceConfigurationResolver) CreateMobileDeviceConfiguration(ctx context.Context, args CreateMobileDeviceConfigurationInput) (*MobileDeviceConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	_, _, err := imeiService.FindById(imeiconfiguration.GetImeiConfigurationData{
		ID:        uuid.MustParse(string(args.ReferenceImeiConfigurationId)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	_, _, err = imsiService.FindById(imsiconfiguration.GetImsiConfigurationByIdData{
		ID:        uuid.MustParse(string(args.ReferenceImsiConfigurationId)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	mobileDeviceConfiguration, _, err := mobileDeviceService.Create(CreateMobileDeviceConfigurationData{
		ProjectId:                    uuid.MustParse(authorization.ProjectId),
		CreatedBy:                    authorization.AccountId,
		Title:                        args.Title,
		Msisdn:                       args.Msisdn,
		ReferenceImsiConfigurationId: uuid.MustParse(string(args.ReferenceImsiConfigurationId)),
		ReferenceImeiConfigurationId: uuid.MustParse(string(args.ReferenceImeiConfigurationId)),
		PermittedLabel:               model.DevicePermittedLabel(args.PermittedLabel),
		BlacklistPriority:            model.BlacklistPriority(args.BlacklistPriority),
		Tags:                         args.Tags,
		StationLocationId:            uuid.MustParse(string(args.StationLocationId)),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return mobileDeviceConfiguration, nil
}

func (parent MobileDeviceConfiguration) ReferenceImeiConfiguration(ctx context.Context) *imeiconfiguration.ImeiConfiguration {
	thunk := imeiDataloader.Load(ctx, dataloader.StringKey(parent.ReferenceImeiConfigurationId))
	imei, err := thunk()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	imeiValue := imei.(imeiconfiguration.ImeiConfiguration)

	return &imeiValue
}

func (parent MobileDeviceConfiguration) ReferenceImsiConfiguration(ctx context.Context) *imsiconfiguration.Imsiconfiguration {
	thunk := imsiDataloader.Load(ctx, dataloader.StringKey(parent.ReferenceImsiConfigurationId))
	imsi, err := thunk()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	imsiValue := imsi.(imsiconfiguration.Imsiconfiguration)

	return &imsiValue
}

func (parent MobileDeviceConfiguration) Tags(ctx context.Context) (*[]tag.Tag, error) {
	thunk := tagDataloader.Load(ctx, dataloader.StringKey(parent.ID))
	tagsLoaderResult, err := thunk()
	if err != nil {
		log.Println(err.Error())
		return nil, nil
	}

	tags := lo.Map(tagsLoaderResult.([]tag.MobileDeviceTag), func(item tag.MobileDeviceTag, index int) tag.Tag {
		return tag.Tag{
			Id:        graphql.ID(item.Tag.ID.String()),
			Title:     item.Tag.Title,
			ProjectId: graphql.ID(item.Tag.ProjectId.String()),
		}
	})

	return &tags, nil
}
