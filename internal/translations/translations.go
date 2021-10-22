package translations

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	message.SetString(language.English, "publication_sorts.year.desc", "Year (newest first)")
	message.SetString(language.English, "publication_sorts.date_created.desc", "Added (newest first)")
	message.SetString(language.English, "publication_sorts.date_updated.desc", "Updated (newest first)")
}
