package helpers

import (
	"html/template"
	"time"

	"github.com/rvflash/elapsed"
	"github.com/ugent-library/biblio-backend/internal/engine"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"searchArgs":  engine.NewSearchArgs,
		"timeElapsed": timeElapsed,
	}
}

func timeElapsed(timestamp string) (string, error) {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "", err
	}
	return elapsed.LocalTime(t, "en"), nil
}
