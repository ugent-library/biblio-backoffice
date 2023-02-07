package helpers

import (
	"fmt"
	"html/template"
	"math"

	"github.com/rvflash/elapsed"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"searchArgs":  models.NewSearchArgs,
		"timeElapsed": elapsed.LocalTime,
		"formatRange": FormatRange,
		"formatBool":  FormatBool,
		"formatBytes": FormatBytes,
	}
}

func FormatRange(start, end string) string {
	var v string
	if len(start) > 0 && len(end) > 0 && start == end {
		v = start
	} else if len(start) > 0 && len(end) > 0 {
		v = fmt.Sprintf("%s - %s", start, end)
	} else if len(start) > 0 {
		v = fmt.Sprintf("%s -", start)
	} else if len(end) > 0 {
		v = fmt.Sprintf("- %s", end)
	}

	return v
}

func FormatBool(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// based on https://github.com/dustin/go-humanize/blob/v1.0.0/bytes.go
// but with KB instead of kB + hide .0 remainders
func FormatBytes(n int64) string {
	if n < 10 {
		return fmt.Sprintf("%d B", n)
	}
	e := math.Floor(math.Log(float64(n)) / math.Log(1000))
	unit := byteUnits[int(e)]
	val := math.Floor(float64(n)/math.Pow(1000, e)*10+0.5) / 10
	format := "%.0f %s"
	if val < 10 && val != math.Trunc(val) {
		format = "%.1f %s"
	}

	return fmt.Sprintf(format, val, unit)
}
