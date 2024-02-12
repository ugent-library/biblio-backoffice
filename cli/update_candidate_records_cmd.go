package cli

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/recordsources"
	_ "github.com/ugent-library/biblio-backoffice/recordsources/plato"
)

func init() {
	rootCmd.AddCommand(updateCandidateRecords)
}

var updateCandidateRecords = &cobra.Command{
	Use:   "update-candidate-records",
	Short: "Update candidate records",
	RunE: func(cmd *cobra.Command, args []string) error {
		services := newServices()

		for _, name := range []string{"plato"} {
			src, err := recordsources.New(name, "https://plato.ea.ugent.be/service/dr/2biblio.jsp")
			if err != nil {
				return err
			}

			// TODO make record mappers
			err = src.GetRecords(context.Background(), func(srcRec recordsources.Record) error {
				p := &models.Publication{}
				var assignedUserID string
				md := gjson.ParseBytes(srcRec.SourceMetadata)

				p.Type = "dissertation"
				p.Status = "private"
				p.Classification = "U"

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
				if ugentID == "" && md.Get("student.studid").String() != "" {
					ugentID = "0000" + md.Get("student.studid").String()
				}
				if ugentID != "" {
					hits, err := services.PersonSearchService.SuggestPeople(ugentID)
					if err != nil {
						return err
					}
					if len(hits) != 1 {
						return errors.New("multiple or no matches for ugent id " + ugentID)
					}
					c := models.ContributorFromPerson(hits[0])
					p.Author = append(p.Author, c)
					assignedUserID = c.PersonID
				} else {
					c := models.ContributorFromFirstLastName(md.Get("student.first").String(), md.Get("student.last").String())
					c.ExternalPerson.Affiliation = md.Get("student.affil").String()
					c.ExternalPerson.HonorificPrefix = md.Get("student.title").String()
					p.Author = append(p.Author, c)
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
						c := models.ContributorFromFirstLastName(val.Get("first").String(), val.Get("last").String())
						c.ExternalPerson.Affiliation = val.Get("affil").String()
						c.ExternalPerson.HonorificPrefix = val.Get("title").String()
						p.Supervisor = append(p.Supervisor, c)
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
					p.AddAbstract(&models.Text{Lang: "und", Text: v.String()})
				}
				if v := md.Get("pdf.url"); v.Exists() {
					sha256, size, err := storeFile(context.TODO(), services.FileStore, v.String())
					if err != nil {
						return err
					}
					// TODO visibility
					f := &models.PublicationFile{
						Relation:    "main_file",
						Name:        srcRec.SourceID + ".pdf",
						ContentType: "application/pdf",
						Size:        size,
						SHA256:      sha256,
					}
					embargo := md.Get("pdf.embargo").String()
					access := md.Get("pdf.accesstype").String()
					if strings.HasPrefix(embargo, "9999") {
						f.AccessLevel = "info:eu-repo/semantics/closedAccess"
					} else if embargo != "" {
						f.AccessLevel = "info:eu-repo/semantics/embargoedAccess"
						f.AccessLevelDuringEmbargo = "info:eu-repo/semantics/closedAccess"
						f.EmbargoDate = embargo[:10]
						if access == "U" {
							f.AccessLevelAfterEmbargo = "info:eu-repo/semantics/restrictedAccess"
						} else if access == "W" {
							f.AccessLevelAfterEmbargo = "info:eu-repo/semantics/openAccess"
						}
					} else if access == "U" {
						f.AccessLevel = "info:eu-repo/semantics/restrictedAccess"
					} else if access == "W" {
						f.AccessLevel = "info:eu-repo/semantics/openAccess"
					}
					p.AddFile(f)
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
					AssignedUserID: assignedUserID,
				}); err != nil {
					return err
				}

				platoID := md.Get("plato_id").String()
				logger.Infof("added candidate record %s from source plato", platoID)

				return nil
			})

			if err != nil {
				return err
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
