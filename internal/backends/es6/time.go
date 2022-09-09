package es6

import "time"

//time.RFC3339 does not include milliseconds
const TimeFormatUTC = "2006-01-02T15:04:05.999Z"

//force UTC
//force use of milliseconds
func FormatTimeUTC(t *time.Time) string {
	return t.UTC().Format(TimeFormatUTC)
}
