package util

import (
	"fmt"

	"slices"
)

func ParseBoolean(v any) bool {
	switch b := v.(type) {
	case int32:
		return b == 1
	case int64:
		return b == 1
	case string:
		return b == "true" || b == "1"
	case bool:
		return b
	}
	return false
}

func ParseString(v any) string {
	switch s := v.(type) {
	case int:
		return fmt.Sprintf("%d", s)
	case int32:
		return fmt.Sprintf("%d", s)
	case int64:
		return fmt.Sprintf("%d", s)
	case float32:
		return fmt.Sprintf("%g", s)
	case float64:
		return fmt.Sprintf("%g", s)
	case string:
		return s
	case bool:
		if s {
			return "true"
		} else {
			return "false"
		}
	}
	return ""
}

func UniqStrings(vals []string) (newVals []string) {
	for _, val := range vals {
		if !slices.Contains(newVals, val) {
			newVals = append(newVals, val)
		}
	}
	return
}
