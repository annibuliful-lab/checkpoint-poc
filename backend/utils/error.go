package utils

import "errors"

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	InternalServerError    = errors.New("internal server error")
	SignTokenFailed        = errors.New("sign token is failed")
	InvalidToken           = errors.New("invalid token")
	TokenExpire            = errors.New("token has expire")
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
	ForbiddenOperation     = errors.New("forbidden operation")
	ContactOwner           = errors.New("please contact project owner")
)
