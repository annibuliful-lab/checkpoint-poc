package utils

type OffsetPagination struct {
	Skip  int64 `json:"skip"`
	Limit int64 `json:"limit"`
}
