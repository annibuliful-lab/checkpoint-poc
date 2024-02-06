package projectRole

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"context"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type ProjectRoleService struct{}

func (ProjectRoleService) Delete(data DeleteProjectRoleData) error {
	dbClient := db.GetPrimaryClient()
	now := time.Now()
	softDeleteStmt := table.ProjectRole.
		UPDATE(table.ProjectRole.DeletedAt).
		MODEL(model.ProjectRole{
			DeletedAt: &now,
		}).
		WHERE(table.ProjectRole.ID.EQ(pg.UUID(uuid.MustParse(string(data.ID)))).
			AND(table.ProjectRole.ProjectId.EQ(pg.UUID(data.ProjectId))))

	result, err := softDeleteStmt.Exec(dbClient)

	if err != nil {
		log.Println(err.Error())
		return utils.InternalServerError
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		log.Println(err.Error())
		return utils.InternalServerError
	}

	if rowAffected == 0 {
		return utils.IdNotfound
	}

	return nil
}
func (ProjectRoleService) Create(data CreateProjectRoleData) (*ProjectRole, string, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	insertProjectRoleStmt := table.ProjectRole.
		INSERT(table.ProjectRole.ID, table.ProjectRole.ProjectId, table.ProjectRole.Title).
		MODEL(model.ProjectRole{
			ID:        uuid.New(),
			Title:     data.Title,
			ProjectId: data.ProjectId,
		}).
		RETURNING(table.ProjectRole.AllColumns)

	projectRole := model.ProjectRole{}

	err = insertProjectRoleStmt.QueryContext(ctx, tx, &projectRole)

	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	for _, id := range data.PermissionIds {
		insertRolePermissionStmt := table.ProjectRolePermission.
			INSERT(table.ProjectRolePermission.ID, table.ProjectRolePermission.RoleId, table.ProjectRolePermission.PermissionId, table.ProjectRole.ProjectId).
			MODEL(model.ProjectRolePermission{
				ID:           uuid.New(),
				RoleId:       projectRole.ID,
				PermissionId: uuid.MustParse(string(id)),
				ProjectId:    data.ProjectId,
			})

		_, err = insertRolePermissionStmt.ExecContext(ctx, tx)

		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
			return nil, utils.InternalServerError.Error(), utils.InternalServerError
		}

	}

	tx.Commit()

	return &ProjectRole{
		Id:        graphql.ID(projectRole.ID.String()),
		ProjectId: graphql.ID(projectRole.ProjectId.String()),
		Title:     projectRole.Title,
		CreatedAt: graphql.Time{Time: projectRole.CreatedAt},
	}, "Created", nil
}

func (ProjectRoleService) Update(data UpdateProjectRoleData) (*ProjectRole, string, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	deleteAllRolePermissionStmt := table.ProjectRolePermission.
		DELETE().
		WHERE(table.ProjectRolePermission.RoleId.EQ(pg.UUID(data.ID)))

	_, err = deleteAllRolePermissionStmt.ExecContext(ctx, tx)

	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	updateProjectRoleStmt := table.ProjectRole.
		UPDATE(table.ProjectRole.Title, table.ProjectRole.UpdatedAt).
		MODEL(model.ProjectRole{
			Title:     data.Title,
			UpdatedAt: time.Now(),
		}).
		WHERE(table.ProjectRole.ID.EQ(pg.UUID(data.ID))).
		RETURNING(table.ProjectRole.AllColumns)

	projectRole := model.ProjectRole{}
	err = updateProjectRoleStmt.QueryContext(ctx, tx, &projectRole)

	if err != nil && db.HasNoRow(err) {
		log.Println(err.Error())
		tx.Rollback()
		return nil, utils.IdNotfound.Error(), utils.IdNotfound
	}

	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	for _, id := range data.PermissionIds {
		insertRolePermissionStmt := table.ProjectRolePermission.
			INSERT(table.ProjectRolePermission.ID, table.ProjectRolePermission.RoleId, table.ProjectRolePermission.PermissionId, table.ProjectRole.ProjectId).
			MODEL(model.ProjectRolePermission{
				ID:           uuid.New(),
				RoleId:       projectRole.ID,
				PermissionId: uuid.MustParse(string(id)),
				ProjectId:    data.ProjectId,
			})

		_, err = insertRolePermissionStmt.ExecContext(ctx, tx)

		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
			return nil, utils.InternalServerError.Error(), utils.InternalServerError
		}

	}

	tx.Commit()

	return &ProjectRole{
		Id:        graphql.ID(projectRole.ID.String()),
		ProjectId: graphql.ID(projectRole.ProjectId.String()),
		Title:     projectRole.Title,
		CreatedAt: graphql.Time{Time: projectRole.CreatedAt},
	}, "Updated", nil
}

func (ProjectRoleService) FindById(data GetProjectRoleByIdData) (*ProjectRole, string, error) {
	dbClient := db.GetPrimaryClient()

	selectProjectRoleStmt := pg.
		SELECT(table.ProjectRole.AllColumns).
		FROM(table.ProjectRole).
		WHERE(table.ProjectRole.ID.EQ(pg.UUID(data.ID))).
		LIMIT(1)

	projectRole := model.ProjectRole{}
	err := selectProjectRoleStmt.Query(dbClient, &projectRole)

	if err != nil && db.HasNoRow(err) {
		log.Println("select-project-role-error", err.Error())
		return nil, utils.IdNotfound.Error(), utils.IdNotfound
	}

	if err != nil {
		log.Println("select-project-role-error", err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	return &ProjectRole{
		Id:        graphql.ID(projectRole.ID.String()),
		ProjectId: graphql.ID(projectRole.ProjectId.String()),
		Title:     projectRole.Title,
		CreatedAt: graphql.Time{Time: projectRole.CreatedAt},
	}, "Ok", nil
}

func (ProjectRoleService) FindMany(data GetProjectRolesData) (*[]ProjectRole, int, error) {
	dbClient := db.GetPrimaryClient()

	selectProjectRolesConditions := pg.Bool(true)
	{
		if data.Search != nil {
			selectProjectRolesConditions = selectProjectRolesConditions.AND(table.ProjectRole.Title.LIKE(pg.String(*data.Search)))
		}
		selectProjectRolesConditions = selectProjectRolesConditions.AND(table.ProjectRole.ProjectId.EQ(pg.UUID(data.ProjectId)))
	}

	selectProjectRolesStmt := pg.
		SELECT(table.ProjectRole.AllColumns).
		FROM(table.ProjectRole).
		WHERE(selectProjectRolesConditions).
		LIMIT(data.pagination.Limit).
		OFFSET(data.pagination.Skip).
		ORDER_BY(table.ProjectRole.CreatedAt)

	projectRoles := []model.ProjectRole{}
	err := selectProjectRolesStmt.Query(dbClient, &projectRoles)

	if err != nil && db.HasNoRow(err) {
		return nil, 404, utils.IdNotfound
	}

	if err != nil {
		log.Println("select-project-role-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	projectRolesResponse := lo.Map(projectRoles, func(projectRole model.ProjectRole, index int) ProjectRole {
		return ProjectRole{
			Id:        graphql.ID(projectRole.ID.String()),
			ProjectId: graphql.ID(projectRole.ProjectId.String()),
			Title:     projectRole.Title,
			CreatedAt: graphql.Time{Time: projectRole.CreatedAt},
		}
	})

	return &projectRolesResponse, 200, nil
}

func (ProjectRoleService) GetProjectRolePermissions(id uuid.UUID) (*[]PermissionResponse, string, error) {
	dbClient := db.GetPrimaryClient()

	projectRolePermissionsStmt := pg.
		SELECT(table.Permission.Action, table.Permission.Subject, table.Permission.ID, table.ProjectRolePermission.RoleId).
		FROM(
			table.ProjectRolePermission.
				INNER_JOIN(table.Permission, table.ProjectRolePermission.PermissionId.EQ(table.Permission.ID)),
		).
		WHERE(table.ProjectRolePermission.RoleId.EQ(pg.UUID(id)))

	var permissions []struct {
		RoleId     uuid.UUID `alias:"project_role_permission.roleId"`
		Permission model.Permission
	}

	err := projectRolePermissionsStmt.Query(dbClient, &permissions)

	if err != nil {
		log.Println(err.Error())
		return &[]PermissionResponse{}, utils.InternalServerError.Error(), utils.InternalServerError
	}

	permissionsResponse := lo.Map(permissions, func(item struct {
		RoleId     uuid.UUID `alias:"project_role_permission.roleId"`
		Permission model.Permission
	}, index int) PermissionResponse {
		return PermissionResponse{
			RoleID:  item.RoleId,
			ID:      item.Permission.ID,
			Action:  item.Permission.Action,
			Subject: item.Permission.Subject,
		}
	})

	return &permissionsResponse, "Ok", nil
}

func (ProjectRoleService) GetProjectRolePermissionsByProjectRoleIds(ids []uuid.UUID) (*[]PermissionResponse, int, error) {
	dbClient := db.GetPrimaryClient()

	var projectRoleIds []pg.Expression
	for _, id := range ids {
		projectRoleIds = append(projectRoleIds, pg.UUID(id))
	}

	projectRolePermissionsStmt := pg.
		SELECT(table.Permission.Action, table.Permission.Subject, table.Permission.ID, table.ProjectRolePermission.RoleId).
		FROM(
			table.ProjectRolePermission.
				INNER_JOIN(table.Permission, table.ProjectRolePermission.PermissionId.EQ(table.Permission.ID)),
		).
		WHERE(table.ProjectRolePermission.RoleId.IN(projectRoleIds...))

	var permissions []struct {
		RoleId     uuid.UUID `alias:"project_role_permission.roleId"`
		Permission model.Permission
	}

	err := projectRolePermissionsStmt.Query(dbClient, &permissions)
	if err != nil {
		log.Println(err.Error())
		return &[]PermissionResponse{}, 500, utils.InternalServerError
	}

	permissionsResponse := lo.Map(permissions, func(item struct {
		RoleId     uuid.UUID `alias:"project_role_permission.roleId"`
		Permission model.Permission
	}, index int) PermissionResponse {
		return PermissionResponse{
			RoleID:  item.RoleId,
			ID:      item.Permission.ID,
			Action:  item.Permission.Action,
			Subject: item.Permission.Subject,
		}
	})

	return &permissionsResponse, 200, nil
}
