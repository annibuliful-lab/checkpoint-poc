package imsiconfiguration

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

type ImsiConfigurationResolver struct{}

var imsiConfigurationService = ImsiConfigurationService{}
var tagService = tag.TagService{}
var tagDataloader = tagService.ImsiConfigurationDataloader()

func (ImsiConfigurationResolver) UpsertImsiConfiguration(ctx context.Context, input UpsertImsiConfigurationInput) (*ImsiConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	imsi, err := imsiConfigurationService.Upsert(UpsertImsiConfigurationData{
		UpdatedBy:         authorization.AccountId,
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		StationLocationId: uuid.MustParse(authorization.StationId),
		Imsi:              input.Imsi,
		BlacklistPriority: &input.BlacklistPriority,
		PermittedLabel:    &input.PermittedLabel,
		Tags:              input.Tags,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return imsi, nil
}

func (ImsiConfigurationResolver) DeleteImsiConfiguration(ctx context.Context, args DeleteImsiConfigurationInput) (*utils.DeleteOperation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	_, err := imsiConfigurationService.Delete(DeleteImsiConfigurationData{
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

func (ImsiConfigurationResolver) GetImsiConfigurationById(ctx context.Context, args GetImsiConfigurationByIdInput) (*ImsiConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	imsiConfiguration, _, err := imsiConfigurationService.FindById(GetImsiConfigurationByIdData{
		ID:        uuid.MustParse(string(args.ID)),
		ProjectId: uuid.MustParse(string(authorization.ProjectId)),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return imsiConfiguration, nil
}

func (ImsiConfigurationResolver) CreateImsiConfiguration(ctx context.Context, args CreateImsiConfigurationInput) (*ImsiConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	imsiConfiguration, _, err := imsiConfigurationService.Create(CreateImsiConfigurationData{
		Imsi:              args.Imsi,
		PermittedLabel:    args.PermittedLabel.String(),
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		StationLocationId: uuid.MustParse(string(args.StationLocationId)),
		Tags:              args.Tags,
		CreatedBy:         authorization.AccountId,
		BlacklistPriority: args.BlacklistPriority.String(),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return imsiConfiguration, nil

}

func (ImsiConfigurationResolver) GetImsiConfigurations(ctx context.Context, args GetImsiConfigurationsInput) ([]ImsiConfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	imsiConfigurations, _, err := imsiConfigurationService.FindMany(GetImsiConfigurationsData{
		StationLocationId: uuid.MustParse(string(args.StationLocationId)),
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		Tags:              args.Tags,
		Search:            args.Search,
		Mnc:               args.Mnc,
		Mcc:               args.Mcc,
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

	return imsiConfigurations, nil
}

func (ImsiConfigurationResolver) UpdateImsiConfiguration(ctx context.Context, args UpdateImsiConfigurationInput) (*ImsiConfiguration, error) {
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

	imsiConfiguration, _, err := imsiConfigurationService.Update(UpdateImsiConfigurationData{
		ID:                uuid.MustParse(string(args.ID)),
		Imsi:              &args.Imsi,
		UpdatedBy:         authorization.AccountId,
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		PermittedLabel:    permittedLabel,
		BlacklistPriority: blacklistPriority,
		Tags:              args.Tags,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return imsiConfiguration, nil
}

func (parent ImsiConfiguration) Tags(ctx context.Context) (*[]tag.Tag, error) {
	thunk := tagDataloader.Load(ctx, dataloader.StringKey(parent.ID))
	tagsLoaderResult, err := thunk()
	if err != nil {
		log.Println(err.Error())
		return nil, nil
	}
	tags := lo.Map(tagsLoaderResult.([]tag.ImsiTag), func(item tag.ImsiTag, index int) tag.Tag {
		return tag.Tag{
			Id:        graphql.ID(item.Tag.ID.String()),
			Title:     item.Tag.Title,
			ProjectId: graphql.ID(item.Tag.ProjectId.String()),
		}
	})

	return &tags, nil
}
