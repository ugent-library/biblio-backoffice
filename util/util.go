package util

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/ugent-library/biblio-backoffice/vocabularies"
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

func IsDate(val string) bool {
	_, e := time.Parse("2006-01-02", val)
	return e == nil
}

func IsYear(val string) bool {
	y, e := strconv.Atoi(val)
	return e == nil && y > 0 && len(val) == 4
}

func IsPublicationType(val string) bool {
	return slices.Contains(vocabularies.Map["publication_types"], val)
}

func IsStatus(val string) bool {
	return slices.Contains(vocabularies.Map["publication_statuses"], val)
}

func IsDatasetIdentifierType(val string) bool {
	return slices.Contains(vocabularies.Map["dataset_identifier_types"], val)
}

func IsDatasetAccessLevel(val string) bool {
	return slices.Contains(vocabularies.Map["dataset_access_levels"], val)
}

func IsDatasetLicense(val string) bool {
	return slices.Contains(vocabularies.Map["dataset_licenses"], val)
}
