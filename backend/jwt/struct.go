package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SignedTokenParams struct {
	AccountId string `json:"userId"`
}

type JwtPayload struct {
	AccountId uuid.UUID `json:"userId"`
	ExpiresAt int64     `json:"exp"`
	jwt.RegisteredClaims
}
