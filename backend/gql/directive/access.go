package directive

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/auth"
	"checkpoint/gql/enum"
	"checkpoint/utils"
	"context"
)

type AccessDirective struct {
	RequiredProjectId *bool
	Subject           *string
	Action            *enum.PermissionAction
}

func (h *AccessDirective) ImplementsDirective() string {
	return "access"
}

func (h *AccessDirective) Validate(ctx context.Context, _ interface{}) error {
	authorization := auth.GetAuthorizationContext(ctx)

	if h.RequiredProjectId != nil && authorization.ProjectId == "" {
		return utils.GraphqlError{
			Message: "Project id is required",
		}
	}

	requiredAllFields := h.RequiredProjectId != nil && h.Action != nil && h.Subject != nil

	if !requiredAllFields {
		err := auth.VerifyAuthentication(authorization)
		return err
	}

	if requiredAllFields {
		err := auth.VerifyAuthorization(ctx, authorization, utils.AuthorizationPermissionData{
			PermissionSubject: *h.Subject,
			PermissionAction:  model.PermissionAction(h.Action.String()),
		})

		return err
	}

	return nil
}
