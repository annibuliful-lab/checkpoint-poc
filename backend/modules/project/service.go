package project

import (
	"checkpoint/.gen/checkpoint/public/model"
	table "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type ProjectService struct{}

func (ProjectService) Create(data CreateProjectInput) (*Project, string, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	project := model.Project{}

	insertProjectStmt := table.Project.
		INSERT(table.Project.ID, table.Project.Title, table.Project.CreatedBy).
		MODEL(model.Project{
			Title:     data.Title,
			ID:        uuid.New(),
			CreatedBy: data.AccountId,
		}).
		RETURNING(table.Project.AllColumns)

	err = insertProjectStmt.QueryContext(ctx, tx, &project)

	if err != nil && strings.Contains(err.Error(), "duplicate") {
		return nil, utils.DataConflict.Error(), errors.New("duplicated project title")
	}

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	projectRole := model.ProjectRole{}

	insertProjectRoleStmt := table.ProjectRole.
		INSERT(table.ProjectRole.ID, table.ProjectRole.Title, table.ProjectRole.ProjectId).
		MODEL(model.ProjectRole{
			ID:        uuid.New(),
			Title:     "Owner",
			ProjectId: project.ID,
		}).
		RETURNING(table.ProjectRole.AllColumns)

	err = insertProjectRoleStmt.QueryContext(ctx, tx, &projectRole)

	if err != nil && strings.Contains(err.Error(), "duplicate") {
		tx.Rollback()
		log.Println(err.Error())
		return nil, utils.DataConflict.Error(), errors.New("duplicate role title")
	}

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	insertProjectAccountStmt := table.ProjectAccount.
		INSERT(table.ProjectAccount.ID, table.ProjectAccount.RoleId, table.ProjectAccount.AccountId, table.ProjectAccount.ProjectId, table.ProjectAccount.CreatedBy).
		MODEL(model.ProjectAccount{
			ID:        uuid.New(),
			RoleId:    projectRole.ID,
			AccountId: uuid.MustParse(data.AccountId),
			ProjectId: project.ID,
			CreatedBy: data.AccountId,
		})

	_, err = insertProjectAccountStmt.ExecContext(ctx, tx)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	tx.Commit()

	return &Project{
		ID:        graphql.ID(project.ID.String()),
		Title:     project.Title,
		CreatedAt: graphql.Time{Time: project.CreatedAt},
		CreatedBy: graphql.ID(project.CreatedBy),
	}, "Created", nil
}

func (ProjectService) Update(data UpdateProjectData) (*Project, string, error) {
	dbClient := db.GetPrimaryClient()

	selectDuplcatedTitleStmt := pg.
		SELECT(table.Project.Title).
		FROM(table.Project).
		WHERE(table.Project.Title.EQ(pg.String(data.Title))).
		LIMIT(1)

	rows, err := selectDuplcatedTitleStmt.Exec(dbClient)

	if err != nil {
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	rowsAffected, _ := rows.RowsAffected()

	if rowsAffected != 0 {
		return nil, utils.DataConflict.Error(), errors.New("duplicated title")
	}

	now := time.Now()

	updateProjectStmt := table.Project.
		UPDATE(table.Project.Title, table.Project.UpdatedAt, table.Project.UpdatedBy).
		MODEL(model.Project{
			Title:     data.Title,
			UpdatedAt: &now,
			UpdatedBy: &data.AccountId,
		}).
		WHERE(table.Project.ID.EQ(pg.UUID(data.ID))).
		RETURNING(table.Project.AllColumns)

	project := model.Project{}

	err = updateProjectStmt.Query(dbClient, &project)
	if err != nil {
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	updatedBy := graphql.ID(*project.UpdatedBy)

	return &Project{
		ID:        graphql.ID(project.ID.String()),
		Title:     project.Title,
		CreatedAt: graphql.Time{Time: project.CreatedAt},
		CreatedBy: graphql.ID(project.CreatedBy),
		UpdatedBy: &graphql.NullID{Value: &updatedBy},
		UpdatedAt: &graphql.NullTime{Value: &graphql.Time{Time: *project.UpdatedAt}},
	}, "Updated", nil

}

func (ProjectService) GetById(data GetProjectInput) (*Project, string, error) {

	dbClient := db.GetPrimaryClient()
	selectProjectStmt := pg.
		SELECT(table.Project.AllColumns).
		FROM(table.Project).
		WHERE(table.Project.ID.EQ(pg.UUID(uuid.MustParse(string(data.ID)))).AND(table.Project.DeletedAt.IS_NULL())).
		LIMIT(1)

	project := model.Project{}
	err := selectProjectStmt.Query(dbClient, &project)
	if err != nil {
		log.Println(err.Error())
		return nil, utils.InternalServerError.Error(), utils.InternalServerError
	}

	var updatedBy graphql.ID
	if project.UpdatedAt != nil {
		updatedBy = graphql.ID(*project.UpdatedBy)
	}

	var updatedAt *graphql.NullTime
	if project.UpdatedAt != nil {
		updatedAt = &graphql.NullTime{Value: &graphql.Time{Time: *project.UpdatedAt}}
	}

	return &Project{
		ID:        graphql.ID(project.ID.String()),
		Title:     project.Title,
		CreatedAt: graphql.Time{Time: project.CreatedAt},
		CreatedBy: graphql.ID(project.CreatedBy),
		UpdatedBy: &graphql.NullID{Value: &updatedBy},
		UpdatedAt: updatedAt,
	}, "Ok", nil
}

func (ProjectService) Delete(data DeleteProjectInput) (string, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()

	deleteProjectStmt := table.Project.
		UPDATE(table.Project.DeletedAt, table.Project.DeletedBy).
		MODEL(model.Project{
			DeletedBy: &data.AccountId,
			DeletedAt: &now,
		}).
		WHERE(table.Project.ID.EQ(pg.UUID(uuid.MustParse(string(data.ID)))))

	_, err := deleteProjectStmt.Exec(dbClient)

	if err != nil {
		return utils.InternalServerError.Error(), utils.InternalServerError
	}

	return "Success", nil
}
