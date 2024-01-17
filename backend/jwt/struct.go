package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type SignedTokenParams struct {
	AccountId string `json:"userId"`
}

type JwtPayload struct {
	AccountId string `json:"userId"`
	ExpiresAt int64  `json:"exp"`
	jwt.RegisteredClaims
}
