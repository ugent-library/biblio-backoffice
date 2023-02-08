package validation

import (
	"strconv"
	"time"

	"github.com/ugent-library/biblio-backoffice/internal/vocabularies"
)

func InArray(values []string, val string) bool {
	for _, v := range values {
		if val == v {
			return true
		}
	}
	return false
}

func IsDate(val string) bool {
	_, e := time.Parse("2006-01-02", val)
	return e == nil
}

func IsTime(val string) bool {
	_, e := time.Parse("15:04", val)
	return e == nil
}

func IsYear(val string) bool {
	y, e := strconv.Atoi(val)
	return e == nil && y > 0 && len(val) == 4
}

func IsPublicationType(val string) bool {
	return InArray(vocabularies.Map["publication_types"], val)
}

func IsStatus(val string) bool {
	return InArray(vocabularies.Map["publication_statuses"], val)
}

func IsDatasetAccessLevel(val string) bool {
	return InArray(vocabularies.Map["dataset_access_levels"], val)
}

func IsDatasetLicense(val string) bool {
	return InArray(vocabularies.Map["dataset_licenses"], val)
}
