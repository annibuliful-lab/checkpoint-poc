package projectRole

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"context"
	"fmt"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/samber/lo"
)

func CreateProjectRole(data CreateProjectRoleData) (*ProjectRoleResponse, int, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println(err.Error())
		return nil, iris.StatusInternalServerError, utils.InternalServerError
	}

	insertProjectRoleStmt := ProjectRole.
		INSERT(ProjectRole.ID, ProjectRole.ProjectId, ProjectRole.Title).
		MODEL(model.ProjectRole{
			ID:        uuid.New(),
			Title:     data.Title,
			ProjectId: data.ProjectId,
		}).
		RETURNING(ProjectRole.AllColumns)

	projectRole := model.ProjectRole{}

	err = insertProjectRoleStmt.QueryContext(ctx, tx, &projectRole)

	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return nil, iris.StatusInternalServerError, utils.InternalServerError
	}

	for _, id := range data.PermissionIds {
		insertRolePermissionStmt := ProjectRolePermission.
			INSERT(ProjectRolePermission.ID, ProjectRolePermission.RoleId, ProjectRolePermission.PermissionId, ProjectRole.ProjectId).
			MODEL(model.ProjectRolePermission{
				ID:           uuid.New(),
				RoleId:       projectRole.ID,
				PermissionId: uuid.MustParse(id),
				ProjectId:    data.ProjectId,
			})

		_, err = insertRolePermissionStmt.ExecContext(ctx, tx)

		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
			return nil, iris.StatusInternalServerError, utils.InternalServerError
		}

	}

	tx.Commit()

	return &ProjectRoleResponse{
		ID:        projectRole.ID,
		ProjectId: projectRole.ProjectId,
		Title:     projectRole.Title,
		CreatedAt: projectRole.CreatedAt,
		UpdatedAt: &projectRole.UpdatedAt,
	}, 201, nil
}

func UpdateProjectRole(data UpdateProjectRoleData) (*ProjectRoleResponse, int, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println(err.Error())
		return nil, iris.StatusInternalServerError, utils.InternalServerError
	}

	deleteAllRolePermissionStmt := ProjectRolePermission.
		DELETE().
		WHERE(ProjectRolePermission.RoleId.EQ(pg.UUID(data.ID)))

	_, err = deleteAllRolePermissionStmt.ExecContext(ctx, tx)

	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return nil, iris.StatusInternalServerError, utils.InternalServerError
	}

	updateProjectRoleStmt := ProjectRole.
		UPDATE(ProjectRole.Title, ProjectRole.UpdatedAt).
		MODEL(model.ProjectRole{
			Title:     data.Title,
			UpdatedAt: time.Now(),
		}).
		WHERE(ProjectRole.ID.EQ(pg.UUID(data.ID))).
		RETURNING(ProjectRole.AllColumns)

	projectRole := model.ProjectRole{}
	err = updateProjectRoleStmt.QueryContext(ctx, tx, &projectRole)

	if err != nil && db.HasNoRow(err) {
		log.Println(err.Error())
		tx.Rollback()
		return nil, 404, utils.IdNotfound
	}

	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return nil, iris.StatusInternalServerError, utils.InternalServerError
	}

	for _, id := range data.PermissionIds {
		insertRolePermissionStmt := ProjectRolePermission.
			INSERT(ProjectRolePermission.ID, ProjectRolePermission.RoleId, ProjectRolePermission.PermissionId, ProjectRole.ProjectId).
			MODEL(model.ProjectRolePermission{
				ID:           uuid.New(),
				RoleId:       projectRole.ID,
				PermissionId: uuid.MustParse(id),
				ProjectId:    data.ProjectId,
			})

		_, err = insertRolePermissionStmt.ExecContext(ctx, tx)

		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
			return nil, iris.StatusInternalServerError, utils.InternalServerError
		}

	}

	tx.Commit()

	return &ProjectRoleResponse{
		ID:        projectRole.ID,
		ProjectId: projectRole.ProjectId,
		Title:     projectRole.Title,
		CreatedAt: projectRole.CreatedAt,
		UpdatedAt: &projectRole.UpdatedAt,
	}, 200, nil
}

func GetProjectRoleById(data GetProjectRoleByIdData) (*ProjectRoleResponse, int, error) {
	dbClient := db.GetPrimaryClient()

	selectProjectRoleStmt := pg.
		SELECT(ProjectRole.AllColumns).
		FROM(ProjectRole).
		WHERE(ProjectRole.ID.EQ(pg.UUID(data.ID))).
		LIMIT(1)

	projectRole := model.ProjectRole{}
	err := selectProjectRoleStmt.Query(dbClient, &projectRole)

	if err != nil && db.HasNoRow(err) {
		log.Println("select-project-role-error", err.Error())
		return nil, 404, utils.IdNotfound
	}

	if err != nil {
		log.Println("select-project-role-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	return &ProjectRoleResponse{
		ID:        projectRole.ID,
		ProjectId: projectRole.ProjectId,
		Title:     projectRole.Title,
		CreatedAt: projectRole.CreatedAt,
		UpdatedAt: &projectRole.UpdatedAt,
	}, 200, nil
}

func GetProjectRoles(data GetProjectRolesData) (*[]ProjectRoleResponse, int, error) {
	dbClient := db.GetPrimaryClient()

	selectProjectRolesConditions := pg.Bool(true)
	{
		if data.Search != "" {
			selectProjectRolesConditions = selectProjectRolesConditions.AND(ProjectRole.Title.LIKE(pg.String(data.Search)))
		}
		selectProjectRolesConditions = selectProjectRolesConditions.AND(ProjectRole.ProjectId.EQ(pg.UUID(data.ProjectId)))
	}

	selectProjectRolesStmt := pg.
		SELECT(ProjectRole.AllColumns).
		FROM(ProjectRole).
		WHERE(selectProjectRolesConditions).
		LIMIT(data.pagination.Limit).
		OFFSET(data.pagination.Skip).
		ORDER_BY(ProjectRole.CreatedAt)
	fmt.Println(selectProjectRolesStmt.DebugSql())
	projectRoles := []model.ProjectRole{}
	err := selectProjectRolesStmt.Query(dbClient, &projectRoles)

	if err != nil && db.HasNoRow(err) {
		return nil, 404, utils.IdNotfound
	}

	if err != nil {
		log.Println("select-project-role-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	projectRolesResponse := lo.Map(projectRoles, func(projectRole model.ProjectRole, index int) ProjectRoleResponse {
		return ProjectRoleResponse{
			ID:        projectRole.ID,
			ProjectId: projectRole.ProjectId,
			Title:     projectRole.Title,
			CreatedAt: projectRole.CreatedAt,
			UpdatedAt: &projectRole.UpdatedAt,
		}
	})

	return &projectRolesResponse, 200, nil
}

func GetProjectRolePermissionsByProjectRoleId(id uuid.UUID) (*[]PermissionResponse, int, error) {
	dbClient := db.GetPrimaryClient()

	projectRolePermissionsStmt := pg.
		SELECT(Permission.Action, Permission.Subject, Permission.ID, ProjectRolePermission.RoleId).
		FROM(
			ProjectRolePermission.
				INNER_JOIN(Permission, ProjectRolePermission.PermissionId.EQ(Permission.ID)),
		).
		WHERE(ProjectRolePermission.RoleId.EQ(pg.UUID(id)))

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

func GetProjectRolePermissionsByProjectRoleIds(ids []uuid.UUID) (*[]PermissionResponse, int, error) {
	dbClient := db.GetPrimaryClient()

	var projectRoleIds []pg.Expression
	for _, id := range ids {
		projectRoleIds = append(projectRoleIds, pg.UUID(id))
	}

	projectRolePermissionsStmt := pg.
		SELECT(Permission.Action, Permission.Subject, Permission.ID, ProjectRolePermission.RoleId).
		FROM(
			ProjectRolePermission.
				INNER_JOIN(Permission, ProjectRolePermission.PermissionId.EQ(Permission.ID)),
		).
		WHERE(ProjectRolePermission.RoleId.IN(projectRoleIds...))

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
