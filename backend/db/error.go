package db

import "strings"

func HasNoRow(err error) bool {
	return strings.Contains(err.Error(), "no rows")
}

func IsDuplicate(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}
