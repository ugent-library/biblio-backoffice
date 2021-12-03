package engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-orcid/orcid"
	"golang.org/x/text/language"
)

func (e *Engine) IsOnORCID(orcidID, id string) bool {
	if orcidID == "" {
		return false
	}

	params := url.Values{}
	params.Set("q", fmt.Sprintf(`orcid:"%s" AND handle:"1854/LU-%s"`, orcidID, id))

	res, _, err := e.orcidClient.Search(params)
	if err != nil {
		log.Print(err)
		return false
	}

	return res.NumFound > 0
}

func (e *Engine) AddPublicationToORCID(orcidID, orcidToken string, p *models.Publication) (*models.Publication, error) {
	client := orcid.NewMemberClient(orcid.Config{
		Token:   orcidToken,
		Sandbox: e.Config.ORCIDSandbox,
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

	return e.UpdatePublication(p)
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

	log.Printf("%+v", w)

	return w
}
