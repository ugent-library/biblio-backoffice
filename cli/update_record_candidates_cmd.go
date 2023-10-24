package cli

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/recordsources"
	_ "github.com/ugent-library/biblio-backoffice/recordsources/plato"
)

func init() {
	rootCmd.AddCommand(updateRecordCandidates)
}

var updateRecordCandidates = &cobra.Command{
	Use:   "update-record-candidates",
	Short: "Update record candidates",
	RunE: func(cmd *cobra.Command, args []string) error {
		services := newServices()

		for _, name := range []string{"plato"} {
			src, err := recordsources.New(name, "https://plato.ea.ugent.be/service/dr/2biblio.jsp")
			if err != nil {
				return err
			}
			srcRecs, err := src.GetRecords(context.Background())
			if err != nil {
				return err
			}

			// TODO make record mappers
			for _, srcRec := range srcRecs {
				p := &models.Publication{}
				md := gjson.ParseBytes(srcRec.SourceMetadata)

				if v := md.Get("titel.eng"); v.Exists() {
					p.Title = v.String()
					if v := md.Get("titel.ned"); v.Exists() {
						p.AlternativeTitle = append(p.AlternativeTitle, v.String())
					}
				} else if v := md.Get("titel.ned"); v.Exists() {
					p.Title = v.String()
				}
				if v := md.Get("year"); v.Exists() {
					p.Year = v.String()
				}
				if v := md.Get("defense.date"); v.Exists() {
					p.DefenseDate = v.String()
				}
				p.DefensePlace = "Ghent, Belgium" // TODO

				ugentID := md.Get("student.ugentid").String()
				if ugentID == "" {
					ugentID = md.Get("student.studid").String()
				}
				if ugentID != "" {
					hits, err := services.PersonSearchService.SuggestPeople(ugentID)
					if err != nil {
						return err
					}
					if len(hits) != 1 {
						return errors.New("multiple or no matches for ugent id " + ugentID)
					}
					p.Author = append(p.Author, models.ContributorFromPerson(hits[0]))
				} else {
					fn := md.Get("student.first").String()
					ln := md.Get("student.last").String()
					p.Author = append(p.Author, models.ContributorFromFirstLastName(fn, ln))
				}

				md.Get("supervisors").ForEach(func(key, val gjson.Result) bool {
					if v := val.Get("ugentid"); v.Exists() {
						hits, pErr := services.PersonSearchService.SuggestPeople(v.String())
						if err != nil {
							err = pErr
							return false
						}
						if len(hits) != 1 {
							err = errors.New("multiple or no matches for ugent id " + v.String())
							return false
						}
						p.Supervisor = append(p.Supervisor, models.ContributorFromPerson(hits[0]))
					} else {
						fn := val.Get("first").String()
						ln := val.Get("last").String()
						p.Supervisor = append(p.Supervisor, models.ContributorFromFirstLastName(fn, ln))
					}
					return true
				})
				if err != nil {
					return err
				}

				if v := md.Get("pdf.ISBN"); v.Exists() {
					p.ISBN = append(p.ISBN, v.String())
				}
				if v := md.Get("pdf.abstract"); v.Exists() {
					p.Abstract = append(p.Abstract, &models.Text{Lang: "und", Text: v.String()})
				}
				if v := md.Get("pdf.url"); v.Exists() {
					sha256, size, err := storeFile(context.TODO(), services.FileStore, v.String())
					if err != nil {
						return err
					}
					f := &models.PublicationFile{
						Relation:    "main_file",
						Name:        srcRec.SourceID + ".pdf",
						ContentType: "application/pdf",
						Size:        size,
						SHA256:      sha256,
					}
					p.File = append(p.File, f)
				}

				j, err := json.Marshal(p)
				if err != nil {
					return err
				}

				if err := services.Repo.AddCandidateRecord(context.TODO(), &models.CandidateRecord{
					SourceName:     srcRec.SourceName,
					SourceID:       srcRec.SourceID,
					SourceMetadata: srcRec.SourceMetadata,
					Type:           "Publication",
					Metadata:       j,
				}); err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func storeFile(ctx context.Context, f backends.FileStore, url string) (string, int, error) {
	// TODO timeouts
	res, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()
	cr := &countingReader{r: res.Body}

	sha256, err := f.Add(ctx, cr, "")

	return sha256, cr.n, err
}

type countingReader struct {
	r io.Reader
	n int
}

func (r *countingReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.n += n
	return n, err
}
