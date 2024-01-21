package jwt

import (
	"checkpoint/utils"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func SignToken(p SignedTokenParams) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": p.AccountId,
			"exp":    time.Now().Add(time.Minute * 15).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SignRefreshToken(p SignedTokenParams) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": p.AccountId,
			"exp":    time.Now().Add(time.Hour * 72).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*JwtPayload, error) {
	// Parse token with custom claims
	token, err := jwt.ParseWithClaims(tokenString, &JwtPayload{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		log.Println(err.Error())
		return nil, utils.InternalServerError
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, utils.InvalidToken
	}

	// Extract custom claims
	claims, ok := token.Claims.(*JwtPayload)

	if !ok {
		return nil, utils.InvalidToken
	}

	expirationTime := time.Unix(claims.ExpiresAt, 0)

	// Compare the expiration time with the current time
	if expirationTime.Before(time.Now()) {
		return nil, utils.TokenExpire
	}

	return claims, nil
}
