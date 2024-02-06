package auth

import (
	table "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/jwt"
	"checkpoint/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	pg "github.com/go-jet/jet/v2/postgres"
)

func VerifyProjectAccount(data VerifyProjectAccountData) bool {
	dbClient := db.GetPrimaryClient()
	selectProjectAdminStmt := pg.
		SELECT(table.ProjectAccount.AccountId).
		FROM(table.ProjectAccount.
			INNER_JOIN(table.ProjectRole, table.ProjectAccount.RoleId.EQ(table.ProjectRole.ID))).
		WHERE(table.ProjectAccount.AccountId.EQ(pg.UUID(data.AccountId)).
			AND(table.ProjectAccount.ProjectId.EQ(pg.UUID(data.ID))))

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
		SELECT(table.ProjectAccount.AccountId).
		FROM(table.ProjectAccount.
			INNER_JOIN(table.ProjectRole, table.ProjectAccount.RoleId.EQ(table.ProjectRole.ID))).
		WHERE(table.ProjectAccount.AccountId.EQ(pg.UUID(data.AccountId)).
			AND(table.ProjectAccount.ProjectId.EQ(pg.UUID(data.ID))).
			AND(table.ProjectRole.Title.EQ(pg.String("Owner"))))

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
		SELECT(table.ProjectAccount.AccountId).
		FROM(table.ProjectAccount.
			INNER_JOIN(table.ProjectRole, table.ProjectAccount.RoleId.EQ(table.ProjectRole.ID))).
		WHERE(table.ProjectAccount.AccountId.EQ(pg.UUID(data.AccountId)).
			AND(table.ProjectAccount.ProjectId.EQ(pg.UUID(data.ID))).
			AND(table.ProjectRole.Title.EQ(pg.String("Owner"))))

	result, err := selectProjectRoleStmt.Exec(dbClient)

	if err != nil {
		log.Println("verify-project-role", err.Error())
		return false
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected != 0
}

func AuthContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := GetAuthenticationHeaders(r.Header)
		if headers.Authorization == "" {
			next.ServeHTTP(w, r)
			return
		}

		payload, err := jwt.VerifyToken(headers.Token)
		if err != nil {
			if err.Error() == utils.TokenExpire.Error() {
				// Token expired error
				errorResponse := map[string]string{"error": "Token expired"}
				writeJSONResponse(w, errorResponse, http.StatusUnauthorized)
				return
			}
			// Other token verification errors
			errorResponse := map[string]string{"error": "Token verification failed"}
			writeJSONResponse(w, errorResponse, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "token", headers.Token)
		ctx = context.WithValue(ctx, "projectId", headers.ProjectId)
		ctx = context.WithValue(ctx, "accountId", payload.AccountId.String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAuthenticationHeaders(header http.Header) AuthenticationHeader {
	authorization := header.Get("authorization")
	projectId := header.Get("x-project-id")

	return AuthenticationHeader{
		Authorization: authorization,
		ProjectId:     projectId,
		Token:         GetAuthToken(authorization),
	}
}

func GetAuthToken(authorization string) string {
	return strings.Replace(authorization, "Bearer ", "", 1)
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
