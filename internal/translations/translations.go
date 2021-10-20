package translations

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	message.SetString(language.English, "year.desc", "Year (newest first)")
	message.SetString(language.English, "date_created.desc", "Added (newest first)")
	message.SetString(language.English, "date_updated.desc", "Updated (newest first)")
}
