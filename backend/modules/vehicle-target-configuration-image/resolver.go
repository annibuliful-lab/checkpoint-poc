package vehicletargetconfigurationimage

import (
	"checkpoint/auth"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
)

type VehicleTargetConfigurationImageResolver struct{}

var vehicleTargetImageService = VehicleTargetConfigurationImageService{}

func (VehicleTargetConfigurationImageResolver) DeleteVehicleTargetConfigurationImage(ctx context.Context, input DeleteVehicleTargetConfigurationImageInput) (*utils.DeleteOperation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	_, err := vehicleTargetImageService.Delete(ctx, DeleteVehicleTargetConfigurationImageData{
		Id:        uuid.MustParse(string(input.Id)),
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
