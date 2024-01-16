package jwt

import "github.com/golang-jwt/jwt/v5"

type SignedTokenParams struct {
	UserId string `json:"userId"`
}

type JwtPayload struct {
	UserId    string `json:"userId"`
	ExpiresAt int64  `json:"exp"`
	jwt.RegisteredClaims
}
