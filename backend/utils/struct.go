package utils

type OffsetPagination struct {
	Skip  int64 `json:"skip"`
	Limit int64 `json:"limit"`
}

type DeleteOperation struct {
	Success bool `json:"success"`
}
