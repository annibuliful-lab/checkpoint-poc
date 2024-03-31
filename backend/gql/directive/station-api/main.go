package directive

import (
	"checkpoint/auth"
	"context"
	"errors"

	"github.com/graph-gophers/graphql-go/directives"
)

type StationApiDirective struct{}

func (h *StationApiDirective) ImplementsDirective() string {
	return "stationApiAccess"
}

func (h *StationApiDirective) Resolve(ctx context.Context, args interface{}, next directives.Resolver) (output interface{}, err error) {
	authentication := auth.GetStationAuthContext(ctx)

	if authentication.DeviceId == "" {
		return nil, errors.New("device id is required")
	}

	if authentication.ApiKey == "" {
		return nil, errors.New("api key is required")
	}

	stationLocation, err := auth.VerifyStationApiAuthentication(authentication.ApiKey)

	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, "stationId", stationLocation.StationId)
	ctx = context.WithValue(ctx, "projectId", stationLocation.ProjectId)

	return next.Resolve(ctx, args)
}
