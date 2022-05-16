package engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/tasks"
	"github.com/ugent-library/go-orcid/orcid"
	"golang.org/x/text/language"
)

// func (e *Engine) IsOnORCID(orcidID, id string) bool {
// 	if orcidID == "" {
// 		return false
// 	}

// 	params := url.Values{}
// 	params.Set("q", fmt.Sprintf(`orcid:"%s" AND handle:"1854/LU-%s"`, orcidID, id))

// 	res, _, err := e.orcidClient.Search(params)
// 	if err != nil {
// 		log.Print(err)
// 		return false
// 	}

// 	return res.NumFound > 0
// }

// TODO move to controller
func (e *Engine) AddPublicationsToORCID(userID string, s *models.SearchArgs) (string, error) {
	user, err := e.GetUser(userID)
	if err != nil {
		return "", err
	}

	taskID := "orcid:" + uuid.NewString()

	e.Tasks.Add(taskID, func(t tasks.Task) error {
		return e.sendPublicationsToORCIDTask(t, userID, user.ORCID, user.ORCIDToken, s)
	})

	return taskID, nil
}

// TODO make workflow
func (e *Engine) AddPublicationToORCID(orcidID, orcidToken string, p *models.Publication) (*models.Publication, error) {
	client := orcid.NewMemberClient(orcid.Config{
		Token:   orcidToken,
		Sandbox: e.ORCIDSandbox,
	})

	work := publicationToORCID(p)
	putCode, res, err := client.AddWork(orcidID, work)
	if err != nil {
		log.Println(err)
		body, _ := ioutil.ReadAll(res.Body)
		log.Print(string(body))
		return p, err
	}

	p.ORCIDWork = append(p.ORCIDWork, models.PublicationORCIDWork{
		ORCID:   orcidID,
		PutCode: putCode,
	})

	if err := e.Store.UpdatePublication(p); err != nil {
		return nil, err
	}

	return p, nil
}

// TODO move to workflows
func (e *Engine) sendPublicationsToORCIDTask(t tasks.Task, userID, orcidID, orcidToken string, searchArgs *models.SearchArgs) error {
	orcidClient := orcid.NewMemberClient(orcid.Config{
		Token:   orcidToken,
		Sandbox: e.ORCIDSandbox,
	})

	var numDone int

	for {
		hits, _ := e.PublicationSearchService.SearchPublications(searchArgs)

		for _, pub := range hits.Hits {
			numDone++

			var done bool
			for _, ow := range pub.ORCIDWork {
				if ow.ORCID == orcidID { // already sent to orcid
					done = true
					break
				}
			}
			if done {
				continue
			}

			work := publicationToORCID(pub)
			putCode, res, err := orcidClient.AddWork(orcidID, work)
			if res.StatusCode == 409 { // duplicate
				continue
			} else if err != nil {
				return err
			}

			pub.ORCIDWork = append(pub.ORCIDWork, models.PublicationORCIDWork{
				ORCID:   orcidID,
				PutCode: putCode,
			})

			if err := e.Store.UpdatePublication(pub); err != nil {
				return err
			}
		}

		t.Progress(numDone, hits.Total)

		if !hits.NextPage() {
			break
		}
		searchArgs.Page = searchArgs.Page + 1
	}

	return nil
}

func publicationToORCID(p *models.Publication) *orcid.Work {
	w := &orcid.Work{
		URL:     orcid.String(fmt.Sprintf("https://biblio.ugent.be/publication/%s", p.ID)),
		Country: orcid.String("BE"),
		ExternalIDs: &orcid.ExternalIDs{
			ExternalID: []orcid.ExternalID{{
				Type:         "handle",
				Relationship: "SELF",
				Value:        fmt.Sprintf("http://hdl.handle.net/1854/LU-%s", p.ID),
			}},
		},
		Title: &orcid.Title{
			Title: orcid.String(p.Title),
		},
		PublicationDate: &orcid.PublicationDate{
			Year: orcid.String(p.Year),
		},
	}

	for _, role := range []string{"author", "editor"} {
		for _, c := range p.Contributors(role) {
			wc := orcid.Contributor{
				CreditName: orcid.String(c.FullName),
				Attributes: &orcid.ContributorAttributes{
					Role: strings.ToUpper(role),
				},
			}
			if c.ORCID != "" {
				wc.ORCID = &orcid.URI{Path: c.ORCID}
			}
			if w.Contributors == nil {
				w.Contributors = &orcid.Contributors{}
			}
			w.Contributors.Contributor = append(w.Contributors.Contributor, wc)
		}
	}

	switch p.Type {
	case "journal_article":
		w.Type = "JOURNAL_ARTICLE"
	case "book":
		w.Type = "BOOK"
	case "book_chapter":
		w.Type = "BOOK_CHAPTER"
	case "book_editor":
		w.Type = "EDITED_BOOK"
	case "dissertation":
		w.Type = "DISSERTATION"
	case "conference":
		switch p.ConferenceType {
		case "meetingAbstract":
			w.Type = "CONFERENCE_ABSTRACT"
		case "poster":
			w.Type = "CONFERENCE_POSTER"
		default:
			w.Type = "CONFERENCE_PAPER"
		}
	case "miscellaneous":
		switch p.MiscellaneousType {
		case "bookReview":
			w.Type = "BOOK_REVIEW"
		case "report":
			w.Type = "REPORT"
		default:
			w.Type = "OTHER"
		}
	default:
		w.Type = "OTHER"
	}

	if len(p.AlternativeTitle) > 0 {
		w.Title.SubTitle = orcid.String(p.AlternativeTitle[0])
	}

	if len(p.Abstract) > 0 {
		w.ShortDescription = p.Abstract[0].Text
	}

	if p.DOI != "" {
		w.ExternalIDs.ExternalID = append(w.ExternalIDs.ExternalID, orcid.ExternalID{
			Type:         "doi",
			Relationship: "SELF",
			Value:        p.DOI,
		})
	}

	if len(p.Language) > 0 {
		if tag, err := language.Parse(p.Language[0]); err == nil {
			w.LanguageCode = tag.String()
		}
	}

	// log.Printf("%+v", w)

	return w
}
