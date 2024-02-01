package auth

import (
	. "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"log"

	pg "github.com/go-jet/jet/v2/postgres"
)

func VerifyProjectAccount(data VerifyProjectAccountData) bool {
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

func VerifyProjectOwner(data VerifyProjectAccountData) bool {
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
		log.Println("verify-project-owner", err.Error())
		return false
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected != 0
}

func VerifyProjectRole(data VerifyProjectRoleData) bool {
	dbClient := db.GetPrimaryClient()

	selectProjectRoleStmt := pg.
		SELECT(ProjectAccount.AccountId).
		FROM(ProjectAccount.
			INNER_JOIN(ProjectRole, ProjectAccount.RoleId.EQ(ProjectRole.ID))).
		WHERE(ProjectAccount.AccountId.EQ(pg.UUID(data.AccountId)).
			AND(ProjectAccount.ProjectId.EQ(pg.UUID(data.ID))).
			AND(ProjectRole.Title.EQ(pg.String("Admin"))))

	result, err := selectProjectRoleStmt.Exec(dbClient)

	if err != nil {
		log.Println("verify-project-role", err.Error())
		return false
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected != 0
}
