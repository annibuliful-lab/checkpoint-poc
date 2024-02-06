package projectRole

import (
	"checkpoint/auth"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type ProjectRoleResolver struct{}

var projectRoleService = ProjectRoleService{}

func (*ProjectRoleResolver) GetProjectRoles(ctx context.Context, args GetProjectRolesInput) ([]ProjectRole, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	projectRoles, _, err := projectRoleService.FindMany(GetProjectRolesData{
		ProjectId: uuid.MustParse(authorization.ProjectId),
		Search:    args.Search,
		pagination: utils.OffsetPagination{
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

	return *projectRoles, nil
}

func (*ProjectRoleResolver) DeleteProjectRole(ctx context.Context, args struct{ Id graphql.ID }) (*utils.DeleteOperation, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(authorization.ProjectId),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}
	err := projectRoleService.Delete(DeleteProjectRoleData{
		ID:        args.Id,
		ProjectId: uuid.MustParse(authorization.ProjectId),
		AccountId: uuid.MustParse(authorization.AccountId),
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

func (*ProjectRoleResolver) GetProjectRoleById(ctx context.Context, args struct{ Id graphql.ID }) (*ProjectRole, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(authorization.ProjectId),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	projectRole, code, err := projectRoleService.FindById(GetProjectRoleByIdData{
		ID:        uuid.MustParse(string(args.Id)),
		ProjectId: uuid.MustParse(authorization.ProjectId),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    code,
			Message: err.Error(),
		}
	}

	return projectRole, nil
}

func (*ProjectRoleResolver) CreateProjectRole(ctx context.Context, args CreateProjectRoleInput) (*ProjectRole, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(authorization.ProjectId),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	projectRole, code, err := projectRoleService.Create(CreateProjectRoleData{
		ProjectId:     uuid.MustParse(authorization.ProjectId),
		PermissionIds: args.PermissionIds,
		Title:         args.Title,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    code,
			Message: err.Error(),
		}
	}

	return projectRole, nil
}

func (*ProjectRoleResolver) UpdateProjectRole(ctx context.Context, args UpdateProjectRoleInput) (*ProjectRole, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(authorization.ProjectId),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	projectRole, code, err := projectRoleService.Update(UpdateProjectRoleData{
		ID:            uuid.MustParse(string(args.Id)),
		ProjectId:     uuid.MustParse(authorization.ProjectId),
		PermissionIds: args.PermissionIds,
		Title:         args.Title,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    code,
			Message: err.Error(),
		}
	}

	return projectRole, nil
}

func (r ProjectRole) Permissions() ([]ProjectRolePermission, error) {

	permissions, code, err := projectRoleService.GetProjectRolePermissions(uuid.MustParse(string(r.Id)))

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    code,
			Message: err.Error(),
		}
	}

	projectRolePermissions := lo.Map(*permissions, func(item PermissionResponse, index int) ProjectRolePermission {
		return ProjectRolePermission{
			Id:      graphql.ID(item.ID.String()),
			RoleId:  graphql.ID(item.RoleID.String()),
			Subject: item.Subject,
			Action:  item.Action.String(),
		}
	})

	return projectRolePermissions, nil
}
