package vehicleproperty

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/auth"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
)

var vehiclePropertyService = VehiclePropertyService{}

type VehiclePropertyResolver struct{}

func (VehiclePropertyResolver) GetVehicleProperties(ctx context.Context, input GetVehiclePropertiesInput) ([]*VehicleProperty, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	properties, err := vehiclePropertyService.FindMany(GetVehiclePropertiesData{
		Type:      model.PropertyType(input.Type.String()),
		Search:    input.Search,
		ProjectId: uuid.MustParse(authorization.ProjectId),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return properties, nil
}
