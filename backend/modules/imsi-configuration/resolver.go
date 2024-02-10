package imsiconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/auth"
	"checkpoint/modules/tag"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type ImsiConfigurationResolver struct{}

var imsiConfigurationService = ImsiConfigurationService{}
var tagService = tag.TagService{}

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

func (ImsiConfigurationResolver) GetImsiConfigurationById(ctx context.Context, args GetImsiConfigurationByIdInput) (*Imsiconfiguration, error) {
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

func (ImsiConfigurationResolver) CreateImsiConfiguration(ctx context.Context, args CreateImeiConfigurationInput) (*Imsiconfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	imsiConfiguration, _, err := imsiConfigurationService.Create(CreateImsiConfigurationData{
		Imsi:              args.Imsi,
		PermittedLabel:    args.PermittedLabel,
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		StationLocationId: uuid.MustParse(string(args.StationLocationId)),
		Tags:              args.Tags,
		CreatedBy:         authorization.AccountId,
		Priority:          args.Priority,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return imsiConfiguration, nil

}

func (ImsiConfigurationResolver) GetImsiConfigurations(ctx context.Context, args GetImsiConfigurationsInput) ([]Imsiconfiguration, error) {
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

func (ImsiConfigurationResolver) UpdateImsiConfiguration(ctx context.Context, args UpdateImsiConfigurationInput) (*Imsiconfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	imsiConfiguration, _, err := imsiConfigurationService.Update(UpdateImsiConfigurationData{
		ID:             uuid.MustParse(string(args.ID)),
		Imsi:           &args.Imsi,
		UpdatedBy:      authorization.AccountId,
		ProjectId:      uuid.MustParse(authorization.ProjectId),
		PermittedLabel: args.PermittedLabel,
		Priority:       args.Priority,
		Tags:           args.Tags,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return imsiConfiguration, nil
}

func (parent Imsiconfiguration) Tags() (*[]tag.Tag, error) {
	tagsResponse, err := tagService.FindByImsiConfigurationId(uuid.MustParse(string(parent.ID)))
	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	tags := lo.Map(*tagsResponse, func(item model.Tag, index int) tag.Tag {
		return tag.Tag{
			Id:        graphql.ID(item.ID.String()),
			ProjectId: graphql.ID(item.ProjectId.String()),
			Title:     item.Title,
		}
	})

	return &tags, nil
}
