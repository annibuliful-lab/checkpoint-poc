package project

import (
	"checkpoint/auth"
	"checkpoint/utils"
	"context"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type ProjectResolver struct{}

var projectService = ProjectService{}

func (r *ProjectResolver) CreateProject(ctx context.Context, args struct{ Title string }) (*Project, error) {
	project, code, err := projectService.Create(CreateProjectInput{
		Title:     args.Title,
		AccountId: ctx.Value("accountId").(string),
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    string(code),
			Message: err.Error(),
		}
	}

	return project, nil
}

func (r *ProjectResolver) UpdateProject(ctx context.Context, args UpdateProjectInput) (*Project, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(string(args.Id)),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	project, code, err := projectService.Update(UpdateProjectData{
		ID:    uuid.MustParse(string(args.Id)),
		Title: args.Title,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    code,
			Message: utils.ForbiddenOperation.Error(),
		}
	}
	return project, nil
}

func (r *ProjectResolver) GetProjectById(ctx context.Context, args struct{ Id graphql.ID }) (*Project, error) {
	authorization := auth.GetAuthorizationContext(ctx)

	match := auth.VerifyProjectAccount(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(string(args.Id)),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}
	project, code, err := projectService.GetById(GetProjectInput{
		ID: args.Id,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    code,
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	return project, nil
}

func (r *ProjectResolver) DeleteProject(ctx context.Context, args struct{ Id graphql.ID }) (*utils.DeleteOperation, error) {
	authorization := auth.GetAuthorizationContext(ctx)
	match := auth.VerifyProjectOwner(auth.VerifyProjectAccountData{
		ID:        uuid.MustParse(string(args.Id)),
		AccountId: uuid.MustParse(authorization.AccountId),
	})

	if !match {
		return nil, utils.GraphqlError{
			Code:    utils.ForbiddenOperation.Error(),
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	projectService.Delete(DeleteProjectInput{
		AccountId: authorization.AccountId,
		ID:        args.Id,
	})

	return &utils.DeleteOperation{
		Success: true,
	}, nil

}
