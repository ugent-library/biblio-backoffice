package views

import (
  "html/template"
  "github.com/unrolled/render"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type RenderAble interface {
	Render(render * render.Render) (template.HTML, error)
}

type TextInput struct {
	Name     string
  Value    string
  Label    string
  Required bool
  Tooltip  string
  Cols     int
  HasError bool
  Error    models.FormError
}

type MultiTextInput struct {
	Name     string
  Values   []string
  Label    string
  Required bool
  Tooltip  string
  Cols     int
  HasError bool
  Error    models.FormError
}

type SelectOption struct {
  Value    string
  Label    string
  Selected bool
}

type Select struct {
  Name     string
  Values   []*SelectOption
  Label    string
  Required bool
  Tooltip  string
  Cols     int
  HasError bool
  Error    models.FormError
}

type MultiSelect struct {
  Name       string
  Values     [][]*SelectOption
  Vocabulary []*SelectOption
  Label      string
  Required   bool
  Tooltip    string
  Cols       int
  HasError   bool
  Error      models.FormError
}

type CheckboxInput struct {
	Name 			string
  Value 		string
  Checked 	bool
	Label    	string
  Required 	bool
  Tooltip  	string
  Cols     	int
  HasError 	bool
  Error   	models.FormError
}

func (e * TextInput) Render(render * render.Render) (template.HTML, error) {
	return RenderPartial(render, "form/_text", e)
}

func (e * MultiTextInput) Render(render * render.Render) (template.HTML, error) {
  return RenderPartial(render, "form/_text_multiple", e)
}

func (e * Select) Render(render * render.Render) (template.HTML, error) {
  return RenderPartial(render, "form/_list", e)
}

func (e * MultiSelect) Render(render * render.Render) (template.HTML, error) {
  return RenderPartial(render, "form/_list_multiple", e)
}

func (e * CheckboxInput) Render(render * render.Render) (template.HTML, error) {
  return RenderPartial(render, "form/_checkbox", e)
}
