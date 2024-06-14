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
	formValue := strings.ReplaceAll(strings.TrimSpace(r.FormValue("mutations")), "\r\n", "\n")
	lines := strings.Split(formValue, "\n")

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

LINES:
	for lineIndex, line := range lines {
		line = strings.TrimSpace(line)

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		reader := csv.NewReader(strings.NewReader(line))
		reader.TrimLeadingSpace = true
		rec, err := reader.Read()

		if err != nil {
			errorMsgs = append(errorMsgs, fmt.Sprintf("error parsing line %d", lineIndex+1))
			continue
		}

		if len(rec) < 2 {
			errorMsgs = append(errorMsgs, fmt.Sprintf("error parsing line %d", lineIndex+1))
			continue
		}

		id := strings.TrimSpace(rec[0])
		op := strings.TrimSpace(rec[1])
		args := rec[2:]
		for i, arg := range args {
			args[i] = strings.TrimSpace(arg)
			if args[i] == "" {
				errorMsgs = append(errorMsgs, fmt.Sprintf("argument %d is empty at line %d", i+1, lineIndex+1))
				continue LINES
			}
		}

		if id == "" {
			errorMsgs = append(errorMsgs, fmt.Sprintf("empty id at line %d", lineIndex+1))
			continue
		}

		if currentID != "" && id != currentID {
			var argErr *mutate.ArgumentError
			err := c.Repo.MutatePublication(currentID, c.User, mutations...)
			if err == nil {
				done++
			} else if errors.As(err, &argErr) {
				errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s: %s", currentID, argErr.Msg))
			} else if len(mutations) == 1 {
				c.Log.Error("could not process publication batch", "id", currentID, "error", err)
				errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s at line %d", currentID, mutations[0].Line))
			} else {
				c.Log.Error("could not process publication batch", "id", currentID, "error", err)
				errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s at lines %d-%d", currentID, mutations[0].Line, mutations[len(mutations)-1].Line))
			}
			mutations = nil
		}

		currentID = id
		mutations = append(mutations, repositories.Mutation{
			Name: op,
			Args: args,
			Line: lineIndex + 1,
		})
	}

	if len(mutations) > 0 {
		var argErr *mutate.ArgumentError
		err := c.Repo.MutatePublication(currentID, c.User, mutations...)
		if err == nil {
			done++
		} else if errors.As(err, &argErr) {
			errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s: %s", currentID, argErr.Msg))
		} else if len(mutations) == 1 {
			c.Log.Error("could not process publication batch", "id", currentID, "error", err)
			errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s at line %d", currentID, mutations[0].Line))
		} else {
			c.Log.Error("could not process publication batch", "id", currentID, "error", err)
			errorMsgs = append(errorMsgs, fmt.Sprintf("could not process publication %s at lines %d-%d", currentID, mutations[0].Line, mutations[len(mutations)-1].Line))
		}
	}

	if len(errorMsgs) == 0 {
		formValue = ""
	}

	publicationviews.BatchBody(c, formValue, done, errorMsgs).Render(r.Context(), w)
}
