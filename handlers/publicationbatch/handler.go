package publicationbatch

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/mutate"
	"github.com/ugent-library/biblio-backoffice/repositories"
	publicationviews "github.com/ugent-library/biblio-backoffice/views/publication"
)

func Show(w http.ResponseWriter, r *http.Request) {
	publicationviews.Batch(ctx.Get(r)).Render(r.Context(), w)
}

func Process(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	formValue := r.FormValue("mutations")
	lines := strings.Split(strings.ReplaceAll(formValue, "\r\n", "\n"), "\n")

	if len(lines) > 500 {
		publicationviews.BatchBody(c, formValue, 0, []string{"no more than 500 operations can be processed at a time"})
		return
	}

	var (
		done      int
		errorMsgs []string
		currentID string
		mutations []repositories.Mutation
	)

	for i, line := range lines {
		line = strings.TrimSpace(line)

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		reader := csv.NewReader(strings.NewReader(line))
		rec, err := reader.Read()

		if err != nil {
			errorMsgs = append(errorMsgs, fmt.Sprintf("error parsing line %d", i+1))
			continue
		}

		if len(rec) < 2 {
			errorMsgs = append(errorMsgs, fmt.Sprintf("error parsing line %d", i+1))
			continue
		}

		id := strings.TrimSpace(rec[0])
		op := strings.TrimSpace(rec[1])
		args := rec[2:]
		for i, arg := range args {
			args[i] = strings.TrimSpace(arg)
		}

		if id == "" {
			errorMsgs = append(errorMsgs, fmt.Sprintf("empty id at line %d", i+1))
			continue
		}

		if currentID != "" && id != currentID {
			var argErr *mutate.ArgumentError
			err := c.Repo.MutatePublication(currentID, c.User, mutations...)
			if err == nil {
				done++
			} else if errors.As(err, &argErr) {
				errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s: %s", currentID, argErr.Error()))
			} else if len(mutations) == 1 {
				errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s at line %d", currentID, mutations[0].Line))
			} else {
				errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s at lines %d-%d", currentID, mutations[0].Line, mutations[len(mutations)-1].Line))
			}
			mutations = nil
		}

		currentID = id
		mutations = append(mutations, repositories.Mutation{
			Name: op,
			Args: args,
			Line: i + 1,
		})
	}

	if len(mutations) > 0 {
		var argErr *mutate.ArgumentError
		err := c.Repo.MutatePublication(currentID, c.User, mutations...)
		if err == nil {
			done++
		} else if errors.As(err, &argErr) {
			errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s: %s", currentID, argErr.Error()))
		} else if len(mutations) == 1 {
			errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s at line %d", currentID, mutations[0].Line))
		} else {
			errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s at lines %d-%d", currentID, mutations[0].Line, mutations[len(mutations)-1].Line))
		}
	}

	if len(errorMsgs) == 0 {
		formValue = ""
	}

	publicationviews.BatchBody(c, formValue, done, errorMsgs).Render(r.Context(), w)
}
