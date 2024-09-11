package authority

import (
	"regexp"
)

var (
	regexMultipleSpaces = regexp.MustCompile(`\s+`)
	regexNoBrackets     = regexp.MustCompile(`[\[\]()\{\}]`)
)
