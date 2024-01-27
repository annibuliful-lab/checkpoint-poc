package project

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func CreateProject(data CreateProjectData) (*ProjectResponse, int, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println(err.Error())
		return nil, iris.StatusInternalServerError, utils.InternalServerError
	}

	var project struct {
		model.Project
	}

	insertProjectStmt := Project.
		INSERT(Project.ID, Project.Title).
		MODEL(model.Project{
			Title: data.Title,
			ID:    uuid.New(),
		}).
		RETURNING(Project.AllColumns)

	err = insertProjectStmt.QueryContext(ctx, tx, &project)

	if err != nil && strings.Contains(err.Error(), "duplicate") {
		return nil, iris.StatusConflict, errors.New("duplicated project title")
	}

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return nil, iris.StatusInternalServerError, utils.InternalServerError
	}

	var projectRole struct {
		model.ProjectRole
	}

	insertProjectRoleStmt := ProjectRole.
		INSERT(ProjectRole.ID, ProjectRole.Title, ProjectRole.ProjectId).
		MODEL(model.ProjectRole{
			ID:        uuid.New(),
			Title:     "Admin",
			ProjectId: project.ID,
		}).
		RETURNING(ProjectRole.AllColumns)

	err = insertProjectRoleStmt.QueryContext(ctx, tx, &projectRole)

	if err != nil && strings.Contains(err.Error(), "duplicate") {
		tx.Rollback()
		log.Println(err.Error())
		return nil, iris.StatusConflict, errors.New("duplicate role title")
	}

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return nil, iris.StatusInternalServerError, utils.InternalServerError
	}

	insertProjectAccountStmt := ProjectAccount.
		INSERT(ProjectAccount.ID, ProjectAccount.RoleId, ProjectAccount.AccountId, ProjectAccount.ProjectId).
		MODEL(model.ProjectAccount{
			ID:        uuid.New(),
			RoleId:    projectRole.ID,
			AccountId: uuid.MustParse(data.AccountId),
			ProjectId: project.ID,
		})

	_, err = insertProjectAccountStmt.ExecContext(ctx, tx)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return nil, iris.StatusInternalServerError, utils.InternalServerError
	}

	tx.Commit()

	return &ProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}, iris.StatusCreated, err
}

func UpdateProject(data UpdateProjectData) (*ProjectResponse, int, error) {
	dbClient := db.GetPrimaryClient()

	selectDuplcatedTitleStmt := pg.
		SELECT(Project.Title).
		FROM(Project).
		WHERE(Project.Title.EQ(pg.String(data.Title)))
	rows, err := selectDuplcatedTitleStmt.Exec(dbClient)

	if err != nil {
		log.Println(err.Error())
		return nil, 500, utils.InternalServerError
	}

	rowsAffected, _ := rows.RowsAffected()

	if rowsAffected != 0 {
		return nil, 409, errors.New("duplicated title")
	}

	now := time.Now()

	updateProjectStmt := Project.
		UPDATE(Project.Title, Project.UpdatedAt).
		MODEL(model.Project{Title: data.Title, UpdatedAt: &now}).
		WHERE(Project.ID.EQ(pg.UUID(data.ID))).
		RETURNING(Project.AllColumns)

	var project struct {
		model.Project
	}

	err = updateProjectStmt.Query(dbClient, &project)
	if err != nil {
		return nil, 500, utils.InternalServerError
	}

	return &ProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}, 200, nil

}

func GetProjectById(data GetProjectData) (*ProjectResponse, int, error) {
	dbClient := db.GetPrimaryClient()
	selectProjectStmt := pg.
		SELECT(Project.AllColumns).
		FROM(Project).
		WHERE(Project.ID.EQ(pg.UUID(data.ID))).
		LIMIT(1)

	var project struct {
		model.Project
	}

	err := selectProjectStmt.Query(dbClient, &project)
	if err != nil {
		log.Println(err.Error())
		return nil, 500, utils.InternalServerError
	}

	return &ProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}, 200, nil
}

func DeleteProjectById(data DeleteProjectData) (int, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()

	deleteProjectStmt := Project.
		UPDATE(Project.DeletedAt, Project.DeleteBy).
		MODEL(model.Project{
			DeleteBy:  data.AccountId,
			DeletedAt: &now,
		}).
		WHERE(Project.ID.EQ(pg.UUID(data.ID)))

	_, err := deleteProjectStmt.Exec(dbClient)

	if err != nil {
		return 500, utils.InternalServerError
	}

	return 200, nil
}

func VerifyAccount(data VerifyProjectAccountData) bool {
	dbClient := db.GetPrimaryClient()
	selectProjectAdminStmt := pg.
		SELECT(ProjectAccount.AccountId).
		FROM(ProjectAccount.
			INNER_JOIN(ProjectRole, ProjectAccount.RoleId.EQ(ProjectRole.ID))).
		WHERE(ProjectAccount.AccountId.EQ(pg.UUID(data.AccountId)).
			AND(ProjectAccount.ProjectId.EQ(pg.UUID(data.ID))))

	result, err := selectProjectAdminStmt.Exec(dbClient)

	if err != nil {
		return false
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected != 0
}

func VerifyOwner(data VerifyProjectAccountData) bool {
	dbClient := db.GetPrimaryClient()

	selectProjectAdminStmt := pg.
		SELECT(ProjectAccount.AccountId).
		FROM(ProjectAccount.
			INNER_JOIN(ProjectRole, ProjectAccount.RoleId.EQ(ProjectRole.ID))).
		WHERE(ProjectAccount.AccountId.EQ(pg.UUID(data.AccountId)).
			AND(ProjectAccount.ProjectId.EQ(pg.UUID(data.ID))).
			AND(ProjectRole.Title.EQ(pg.String("Admin"))))

	result, err := selectProjectAdminStmt.Exec(dbClient)

	if err != nil {
		log.Println("verify-owner", err.Error())
		return false
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected != 0
}
