package utils

import (
	"errors"
	"fmt"
)

var (
	InternalServerError    = errors.New("internal server error")
	SignTokenFailed        = errors.New("sign token is failed")
	InvalidToken           = errors.New("invalid token")
	TokenExpire            = errors.New("token has expire")
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
	ForbiddenOperation     = errors.New("forbidden operation")
	ContactOwner           = errors.New("please contact project owner")
	DataConflict           = errors.New("data conflict")
	IdNotfound             = errors.New("id not found")
	Notfound               = errors.New("data not found")
)

type GraphqlError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e GraphqlError) Error() string {
	return fmt.Sprintf("error [%s]: %s", e.Code, e.Message)
}

func (e GraphqlError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	}
}
