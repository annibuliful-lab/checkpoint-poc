package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type GenerateTokenParams struct {
	UserId     string   `json:"userId"`
	ProjectIds []string `json:"projectIds"`
}

func GenerateToken(options GenerateTokenParams) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":     options.UserId,
		"projectIds": options.ProjectIds,
		"nbf":        time.Now().Unix(),                    // Not Before time
		"exp":        time.Now().Add(time.Hour * 1).Unix(), // Token expiration time (1 hour in this example)
	})

	tokenString, err := token.SignedString("hmacSampleSecret")

	if err == nil {
		panic(err.Error())
	}

	return tokenString
}
