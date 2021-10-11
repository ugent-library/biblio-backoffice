package views

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-web/jsonapi"
)

type textData struct {
	Text     string
	Label    string
	Required bool
	Tooltip  string
}

type listData struct {
	List     []string
	Label    string
	Required bool
	Tooltip  string
}

type textFormData struct {
	Key      string
	Text     string
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    jsonapi.Error
}

type textMultipleFormData struct {
	Key      string
	Text     []string
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    jsonapi.Error
}

type listFormValues struct {
	Key      string
	Value    string
	Selected bool
}

type listFormData struct {
	Key      string
	Values   []*listFormValues
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    jsonapi.Error
}

type listMultipleFormData struct {
	Key        string
	Values     map[int][]*listFormValues
	Vocabulary map[string]string
	Label      string
	Required   bool
	Tooltip    string
	Cols       int
	HasError   bool
	Error      jsonapi.Error
}

type boolData struct {
	Value    bool
	Label    string
	Required bool
	Tooltip  string
}

type CheckboxInput struct {
	Name         string
  Value        string
  Checked      bool
	Label        string
  Required     bool
  Tooltip      string
  Cols         int
  HasError     bool
  Error        jsonapi.Error
}

type Data struct {
	User    *models.User
	request *http.Request
}

func NewData(r *http.Request) Data {
	return Data{
		User:    context.User(r.Context()),
		request: r,
	}
}

func (d Data) IsHTMXRequest() bool {
	return d.request.Header.Get("HX-Request") != ""
}
