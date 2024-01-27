package utils

type OffsetPagination struct {
	Skip  *int `json:"skip"`
	Limit *int `json:"limit"`
}
