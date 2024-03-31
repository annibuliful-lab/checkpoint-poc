package imeiconfiguration

import (
	"checkpoint/auth"
	"checkpoint/modules/tag"
	"checkpoint/utils"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type ImeiConfigurationResolver struct{}

var imeiconfigurationService = ImeiConfigurationService{}
var tagService = tag.TagService{}
var tagDataloader = tagService.ImeiConfigurationDataloader()

func (ImeiConfigurationResolver) UpsertImeiConfiguration(ctx context.Context, input UpsertImeiConfigurationInput) (*ImeiConfiguration, error) {
	if !utils.ValidateIMEI(input.Imei) {
		return nil, utils.ErrInvalidIMEI
	}

	authorization := auth.GetAuthorizationContext(ctx)

	imeiConfiguration, err := imeiconfigurationService.Upsert(UpsertImeiConfigurationData{
		UpdatedBy:         authorization.AccountId,
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		StationLocationId: uuid.MustParse(authorization.StationId),
		Imei:              input.Imei,
		PermittedLabel:    input.PermittedLabel,
		BlacklistPriority: input.BlacklistPriority,
		Tags:              input.Tags,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return imeiConfiguration, nil
}

func (ImeiConfigurationResolver) GetImeiConfigurations(ctx context.Context, args GetImeiConfigurationsInput) ([]ImeiConfiguration, error) {

	authorization := auth.GetAuthorizationContext(ctx)

	var permittedLabel *string
	if args.PermittedLabel != nil {
		value := args.PermittedLabel.String()
		permittedLabel = &value
	}

	var blacklistPriority *string
	if args.BlacklistPriority != nil {
		value := args.BlacklistPriority.String()
		blacklistPriority = &value
	}

	imeiConfigurations, _, err := imeiconfigurationService.FindMany(GetImeiConfigurationsData{
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		StationLocationId: uuid.MustParse(string(args.StationLocationId)),
		Tags:              args.Tags,
		Search:            args.Search,
		BlacklistPriority: blacklistPriority,
		PermittedLabel:    permittedLabel,
		Pagination: utils.OffsetPagination{
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

	return imeiConfigurations, nil
}

func (ImeiConfigurationResolver) GetImeiConfigurationById(ctx context.Context, args GetImeiConfigurationInput) (*ImeiConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	imeiConfiguration, _, err := imeiconfigurationService.FindById(GetImeiConfigurationData{
		ID:        uuid.MustParse(string(args.ID)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return imeiConfiguration, nil
}

func (ImeiConfigurationResolver) DeleteImeiConfiguration(ctx context.Context, args DeleteImeiConfigurationInput) (*utils.DeleteOperation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	_, err := imeiconfigurationService.Delete(DeleteImeiConfigurationData{
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

	return &utils.DeleteOperation{
		Success: true,
	}, nil
}

func (ImeiConfigurationResolver) UpdateImeiConfiguration(ctx context.Context, args UpdateImeiConfigurationInput) (*ImeiConfiguration, error) {
	if !utils.ValidateIMEI(args.Imei) {
		return nil, utils.ErrInvalidIMEI
	}

	authorization := auth.GetAuthorizationContext(ctx)
	ImeiConfiguration, _, err := imeiconfigurationService.Update(UpdateImeiConfigurationData{
		ID:                uuid.MustParse(string(args.ID)),
		UpdatedBy:         authorization.AccountId,
		Imei:              args.Imei,
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		BlacklistPriority: args.BlacklistPriority,
		PermittedLabel:    args.PermittedLabel,
		Tags:              args.Tags,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return ImeiConfiguration, nil
}

func (ImeiConfigurationResolver) CreateImeiConfiguration(ctx context.Context, args CreateImeiConfigurationInput) (*ImeiConfiguration, error) {
	if !utils.ValidateIMEI(args.Imei) {
		return nil, utils.ErrInvalidIMEI
	}

	authorization := auth.GetAuthorizationContext(ctx)
	ImeiConfiguration, _, err := imeiconfigurationService.Create(CreateImeiConfigurationData{
		CreatedBy:         authorization.AccountId,
		Imei:              args.Imei,
		StationLocationId: uuid.MustParse(string(args.StationLocationId)),
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		BlacklistPriority: args.BlacklistPriority.String(),
		PermittedLabel:    args.PermittedLabel.String(),
		Tags:              args.Tags,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return ImeiConfiguration, nil
}

func (parent ImeiConfiguration) Tags(ctx context.Context) (*[]tag.Tag, error) {

	thunk := tagDataloader.Load(ctx, dataloader.StringKey(parent.ID))
	tagsLoaderResult, err := thunk()
	if err != nil {
		log.Println(err.Error())
		return nil, nil
	}
	tags := lo.Map(tagsLoaderResult.([]tag.ImeiTag), func(item tag.ImeiTag, index int) tag.Tag {
		return tag.Tag{
			Id:        graphql.ID(item.Tag.ID.String()),
			Title:     item.Tag.Title,
			ProjectId: graphql.ID(item.Tag.ProjectId.String()),
		}
	})

	return &tags, nil
}
