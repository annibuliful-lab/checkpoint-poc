package imsiconfiguration

import (
	"checkpoint/auth"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
)

type ImsiConfigurationResolver struct{}

var imsiConfigurationService = ImsiConfigurationService{}

func (ImsiConfigurationResolver) CreateImsiConfiguration(ctx context.Context, args CreateImeiConfigurationInput) (*Imsiconfiguration, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	imsiConfiguration, _, err := imsiConfigurationService.Create(CreateImsiConfigurationData{
		Imsi:              args.Imsi,
		PermittedLabel:    args.PermittedLabel,
		ProjectId:         uuid.MustParse(authorization.ProjectId),
		StationLocationId: uuid.MustParse(string(args.StationLocationId)),
		Tags:              *args.Tags,
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
