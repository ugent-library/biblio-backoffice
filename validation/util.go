package validation

import (
	"strconv"
	"time"

	"slices"

	"github.com/ugent-library/biblio-backoffice/vocabularies"
)

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
