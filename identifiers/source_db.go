package identifiers

import "log"

type SourceDBType struct{}

func (i *SourceDBType) Validate(sourceDB string) bool {
	return true
}

func (i *SourceDBType) Normalize(sourceDB string) (string, error) {
	return sourceDB, nil
}

func (i *SourceDBType) Resolve(sourceDB string) string {
	switch sourceDB {
	case "plato":
		return "https://plato.ea.ugent.be/"

	default:
		log.Default().Printf("Unknown SourceDB: %s", sourceDB)
		return ""
	}
}
