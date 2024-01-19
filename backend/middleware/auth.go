package middleware

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/jwt"
	"checkpoint/utils"
	"fmt"
	"log"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/kataras/iris/v12"
)

func AuthMiddleware(permission *string) iris.Handler {
	dbClient := db.GetPrimaryClient()
	// redisClient := db.GetRedisClient()

	return func(ctx iris.Context) {
		headers := utils.GetAuthenticationHeaders(ctx)

		if headers.Token == "" {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": utils.InvalidToken,
			})
			return
		}

		payload, err := jwt.VerifyToken(headers.Token)

		if err != nil && err.Error() == utils.TokenExpire {
			updateStmt := SessionToken.
				UPDATE(SessionToken.Revoke).
				MODEL(model.SessionToken{
					Revoke: true,
				}).
				WHERE(SessionToken.Token.EQ(pg.String(headers.Token)))

			fmt.Println(updateStmt.DebugSql())

			_, err := updateStmt.Exec(dbClient)
			if err != nil {
				log.Println(err.Error())
			}

		}

		if err != nil {

			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": err.Error(),
			})
			return
		}
		var account struct {
			model.AccountConfiguration
		}

		selectAccountStmt := pg.
			SELECT(AccountConfiguration.AllColumns).
			FROM(AccountConfiguration).
			WHERE(AccountConfiguration.AccountId.EQ(pg.UUID(payload.AccountId))).
			LIMIT(1)

		err = selectAccountStmt.Query(dbClient, &account)

		if err != nil {
			log.Println(err.Error())
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"message": utils.InternalServerError,
			})
			return
		}

		if !account.IsActive {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{
				"message": "please contact project owner",
			})
			return
		}

		if headers.ProjectId == "" {
			ctx.Next()
			return
		}

	}
}
