package helpers

import "html/template"

// We need this because Unrolled escapes HTML in input.
// see: https://github.com/unrolled/render/issues/41#issuecomment-153567562

func StringToHtml(raw string) template.HTML {
	return template.HTML(raw)
}
