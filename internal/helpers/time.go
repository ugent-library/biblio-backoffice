package helpers

import (
	"time"

	"github.com/rvflash/elapsed"
)

func TimeElapsed(timestamp string) (string, error) {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "", err
	}
	return elapsed.LocalTime(t, "en"), nil
}
