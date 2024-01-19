package utils

import "errors"

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	InternalServerError = "internal server error"
	SignTokenFailed     = "sign token is failed"
	InvalidToken        = "invalid token"
	TokenExpire         = "token has expire"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)
