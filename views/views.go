//go:generate go run github.com/a-h/templ/cmd/templ@v0.2.543 generate
package views

import (
	"encoding/json"
)

func ToJSON(pairs ...string) string {
	data := make(map[string]string)

	for i := 0; i < len(pairs); i += 2 {
		data[pairs[i]] = pairs[i+1]
	}

	bytes, _ := json.Marshal(data)
	return string(bytes)
}
