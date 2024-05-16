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

	var (
		done      int
		errorMsgs []string
		currentID string
		mutations []repositories.Mutation
		lineNum   int
	)

	for i, line := range lines {
		line = strings.TrimSpace(line)

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		lineNum = i + 1

		reader := csv.NewReader(strings.NewReader(line))
		rec, err := reader.Read()

		if err != nil {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>error parsing line %d</p>", lineNum))
			continue
		}

		if len(rec) < 2 {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>error parsing line %d</p>", lineNum))
			continue
		}

		id := strings.TrimSpace(rec[0])
		op := strings.TrimSpace(rec[1])
		args := rec[2:]
		for i, arg := range args {
			args[i] = strings.TrimSpace(arg)
		}

		if id == "" {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>empty id at line %d</p>", lineNum))
			continue
		}

		if currentID != "" && id != currentID {
			err := c.Repo.MutatePublication(currentID, c.User, mutations...)
			if err == nil {
				done++
			} else {
				errorMsgs = append(errorMsgs, fmt.Sprintf("<p>could not process publication %s at line %d</p>", currentID, lineNum-1))
			}
			mutations = nil
		}

		currentID = id
		mutations = append(mutations, repositories.Mutation{
			Op:   op,
			Args: args,
		})
	}

	if len(mutations) > 0 {
		err := c.Repo.MutatePublication(currentID, c.User, mutations...)
		if err == nil {
			done++
		} else {
			errorMsgs = append(errorMsgs, fmt.Sprintf("<p>could not process publication %s at line %d</p>", currentID, lineNum))
		}
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
