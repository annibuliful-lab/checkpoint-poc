package db

import "strings"

func ConvertArrayStringToInput(data []string) *string {

	resultString := "{" + strings.Join(data, ",") + "}"

	return &resultString
}

func ConvertArrayDbStringToArrayString(data *string) []string {
	if data == nil {
		return []string{}
	}

	trimmedString := strings.Trim(*data, "{}")

	return strings.Split(trimmedString, ",")
}
