package publicationbatch

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/render/flash"
	"github.com/ugent-library/biblio-backoffice/repositories"
	"github.com/ugent-library/biblio-backoffice/views"
)

func Show(w http.ResponseWriter, r *http.Request) {
	views.PublicationbatchShow(ctx.Get(r)).Render(r.Context(), w)
}

func Process(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	lines := strings.Split(strings.ReplaceAll(r.FormValue("ops"), "\r\n", "\n"), "\n")

	if len(lines) > 500 {
		c.PersistFlash(w, *flash.SimpleFlash().
			WithLevel("error").
			WithBody("No more than 500 operations can be processed at one time.").
			DismissedAfter(0))
		http.Redirect(w, r, c.PathTo("publication_batch").String(), http.StatusFound)
		return
	}

	done := 0
	var errorMsgs []string

	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		rdr := csv.NewReader(strings.NewReader(strings.TrimSpace(line)))
		rec, err := rdr.Read()

		if err != nil {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>error parsing line %d</p>", i))
			continue
		}

		if len(rec) < 2 {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>error parsing line %d</p>", i))
			continue
		}

		id := strings.TrimSpace(rec[0])
		op := strings.TrimSpace(rec[1])
		args := rec[2:]
		for i, arg := range args {
			args[i] = strings.TrimSpace(arg)
		}

		err = c.Repo.MutatePublication(id, c.User, repositories.Mutation{
			Op:   op,
			Args: args,
		})

		if err != nil {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>could not process publication %s at line %d</p>", id, i))
			continue
		}

		done++
	}

	if done > 0 {
		c.PersistFlash(w, *flash.SimpleFlash().
			WithLevel("success").
			WithBody(template.HTML(fmt.Sprintf("<p>Successfully processed %d publications.</p>", done))))
	}
	if len(errorMsgs) > 0 {
		c.PersistFlash(w, *flash.SimpleFlash().
			WithLevel("error").
			WithBody(template.HTML(strings.Join(errorMsgs, ""))).
			DismissedAfter(0))
	}

	http.Redirect(w, r, c.PathTo("publication_batch").String(), http.StatusFound)
}
