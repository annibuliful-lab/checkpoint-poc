package utils

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	InternalServerError = "internal server error"
	SignTokenFailed     = "sign token is failed"
)
