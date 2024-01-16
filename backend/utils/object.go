package utils

import "strings"

func ConvertJsonObjectKeyToLower(v any) any {
	switch v := v.(type) {
	case []any:
		lv := make([]any, len(v))
		for i := range v {
			lv[i] = ConvertJsonObjectKeyToLower(v[i])
		}
		return lv
	case map[string]any:
		lv := make(map[string]any, len(v))
		for mk, mv := range v {
			lv[strings.ToLower(mk)] = mv
		}
		return lv
	default:
		return v
	}
}
