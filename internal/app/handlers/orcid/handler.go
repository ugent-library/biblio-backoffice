package orcid

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
	"github.com/ugent-library/biblio-backend/internal/tasks"
	"github.com/ugent-library/biblio-backend/internal/ulid"
	"github.com/ugent-library/go-orcid/orcid"
	"golang.org/x/text/language"
)

type Handler struct {
	handlers.BaseHandler
	Tasks                    *tasks.Hub
	Repository               backends.Repository
	PublicationSearchService backends.PublicationSearchService
	Sandbox                  bool
}

type Context struct {
	handlers.BaseContext
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			h.Logger.Warnw("orcid: user is not authorized to access this resource", "user", ctx.User.ID)
			render.Unauthorized(w, r)
			return
		}

		fn(w, r, Context{
			BaseContext: ctx,
		})
	})
}

type BindAdd struct {
	PublicationID string `path:"id"`
}

type YieldAdd struct {
	Context
	Publication *models.Publication
	Flash       flash.Flash
}

type YieldAddAll struct {
	ID      string
	Status  tasks.Status
	Message string
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindAdd{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("add orcid: could not bind request arguments", "errors", err, "request", r, "user", ctx.User.ID)
		render.BadRequest(w, r, err)
		return
	}

	p, err := h.Repository.GetPublication(b.PublicationID)
	if err != nil {
		h.Logger.Errorw("add orcid: could not get publication", "errors", err, "publication", b.PublicationID, "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}
	if !ctx.User.CanViewPublication(p) {
		h.Logger.Warnw("add orcid: user has no permission to view this publication", "publication", b.PublicationID, "user", ctx.User.ID)
		render.Forbidden(w, r)
		return
	}

	var f flash.Flash

	p, err = h.addPublicationToORCID(ctx.User.ORCID, ctx.User.ORCIDToken, p)
	if err != nil {
		if err == orcid.ErrDuplicate {
			h.Logger.Warnw("add orcid: this publicaton is already part of the users orcid works", "publication", b.PublicationID, "user", ctx.User.ID)
			f = flash.Flash{Type: "info", Body: "This publication is already part of your ORCID works."}
		} else {
			h.Logger.Warnw("add orcid: could not add this publication to the users orcid works", "user", "publication", b.PublicationID, "user", ctx.User.ID)
			f = flash.Flash{Type: "error", Body: "Couldn't add this publication to your ORCID works."}
		}
	} else {
		f = flash.Flash{Type: "success", Body: "Successfully added the publication to your ORCID works.",
			DismissAfter: 5000}
	}

	render.View(w, "publication/refresh_orcid_status", YieldAdd{
		Context:     ctx,
		Publication: p,
		Flash:       f,
	})
}

func (h *Handler) AddAll(w http.ResponseWriter, r *http.Request, ctx Context) {
	id, err := h.addPublicationsToORCID(
		ctx.User,
		models.NewSearchArgs().WithFilter("status", "public").WithFilter("author.id", ctx.User.ID),
	)
	if err != nil {
		h.Logger.Errorw("add all orcid: could not add all publications to the users orcid", "user", ctx.User.ID)
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "task/add_status", YieldAddAll{
		ID:      id,
		Status:  tasks.Status{},
		Message: "",
	})
}

// TODO make workflow
// TODO add proper logging once moved to workflows
func (h *Handler) addPublicationToORCID(orcidID, orcidToken string, p *models.Publication) (*models.Publication, error) {
	client := orcid.NewMemberClient(orcid.Config{
		Token:   orcidToken,
		Sandbox: h.Sandbox,
	})

	work := publicationToORCID(p)
	putCode, res, err := client.AddWork(orcidID, work)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		log.Printf("orcid error: %s", body)
		return p, err
	}

	p.ORCIDWork = append(p.ORCIDWork, models.PublicationORCIDWork{
		ORCID:   orcidID,
		PutCode: putCode,
	})

	if err := h.Repository.SavePublication(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (h *Handler) addPublicationsToORCID(user *models.User, s *models.SearchArgs) (string, error) {
	taskID := "orcid:" + ulid.MustGenerate()

	h.Tasks.Add(taskID, func(t tasks.Task) error {
		return h.sendPublicationsToORCIDTask(t, user.ID, user.ORCID, user.ORCIDToken, s)
	})

	return taskID, nil
}

// TODO move to workflows
func (h *Handler) sendPublicationsToORCIDTask(t tasks.Task, userID, orcidID, orcidToken string, searchArgs *models.SearchArgs) error {
	orcidClient := orcid.NewMemberClient(orcid.Config{
		Token:   orcidToken,
		Sandbox: h.Sandbox,
	})

	var numDone int

	for {
		hits, _ := h.PublicationSearchService.Search(searchArgs)

		for _, p := range hits.Hits {
			numDone++

			var done bool
			for _, ow := range p.ORCIDWork {
				if ow.ORCID == orcidID { // already sent to orcid
					done = true
					break
				}
			}
			if done {
				continue
			}

			work := publicationToORCID(p)
			putCode, res, err := orcidClient.AddWork(orcidID, work)
			if res.StatusCode == 409 { // duplicate
				continue
			} else if err != nil {
				body, _ := ioutil.ReadAll(res.Body)
				log.Printf("orcid error: %s", body)
				return err
			}

			p.ORCIDWork = append(p.ORCIDWork, models.PublicationORCIDWork{
				ORCID:   orcidID,
				PutCode: putCode,
			})

			if err := h.Repository.SavePublication(p); err != nil {
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
				CreditName: orcid.String(strings.Join([]string{c.FirstName, c.LastName}, " ")),
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

	return w
}
