package access

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/auth"
	"checkpoint/gql/enum"
	"checkpoint/utils"
	"context"
)

type AccessDirective struct {
	RequiredProjectId *bool
	RequiredStationId *bool
	Subject           *string
	Action            *enum.PermissionAction
}

func (h *AccessDirective) ImplementsDirective() string {
	return "access"
}

func (h *AccessDirective) Validate(ctx context.Context, _ interface{}) error {
	authorization := auth.GetAuthorizationContext(ctx)

	err := auth.VerifyAuthentication(authorization)

	if err != nil {
		return err
	}

	if h.RequiredProjectId != nil && authorization.ProjectId == "" {
		return utils.GraphqlError{
			Message: "Project id is required",
		}
	}

	if h.RequiredStationId != nil && authorization.StationId == "" {
		return utils.GraphqlError{
			Message: "Station id is required",
		}
	}

	if h.RequiredStationId != nil && *h.RequiredStationId {
		err := auth.VerifyProjectStation(ctx, authorization)
		if err != nil {
			return err
		}
	}

	requiredPermission := h.Action != nil && h.Subject != nil

	if requiredPermission {
		err := auth.VerifyAuthorization(ctx, authorization, utils.AuthorizationPermissionData{
			PermissionSubject: *h.Subject,
			PermissionAction:  model.PermissionAction(h.Action.String()),
		})

		return err
	}

	return nil
}
