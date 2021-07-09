package helpers

import (
	"html/template"
	"time"

	"github.com/rvflash/elapsed"
)

func Time() template.FuncMap {
	return template.FuncMap{
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
