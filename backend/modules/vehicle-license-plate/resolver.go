package vehiclelicenseplate

import "context"

type VehiclelicensePlateResolver struct{}

func (VehiclelicensePlateResolver) CreateVehicleLicensePlate(ctx context.Context, input CreateVehicleLicensePlateInput) (*VehicleLicensePlate, error) {
	return &VehicleLicensePlate{}, nil
}
