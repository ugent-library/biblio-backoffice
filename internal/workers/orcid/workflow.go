package orcid

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-orcid/orcid"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/text/language"
)

type Activities struct {
	PublicationSearchService backends.PublicationSearchService
	PublicationService       backends.PublicationService
	OrcidSandbox             bool
}

type Args struct {
	UserID     string
	ORCID      string
	ORCIDToken string
	SearchArgs models.SearchArgs
}

func SendPublicationsToORCIDWorkflow(ctx workflow.Context, args Args) (err error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		HeartbeatTimeout:    10 * time.Second, // needs to be set for heartbeat based progress to work well
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1.0,
			MaximumInterval:    time.Second * 10,
			MaximumAttempts:    3,
		},
	})

	var a *Activities

	if err = workflow.ExecuteActivity(ctx, a.SendPublicationsToORCID, args).Get(ctx, nil); err != nil {
		return err
	}

	return
}

func (a *Activities) SendPublicationsToORCID(ctx context.Context, args Args) error {
	orcidClient := orcid.NewMemberClient(orcid.Config{
		Token:   args.ORCIDToken,
		Sandbox: a.OrcidSandbox,
	})

	searchArgs := args.SearchArgs

	var numDone int

	for {
		hits, _ := a.PublicationSearchService.UserPublications(args.UserID, &searchArgs)

		for _, pub := range hits.Hits {
			numDone++

			var done bool
			for _, ow := range pub.ORCIDWork {
				if ow.ORCID == args.ORCID { // already sent to orcid
					done = true
					break
				}
			}
			if done {
				continue
			}

			work := publicationToORCID(pub)
			putCode, res, err := orcidClient.AddWork(args.ORCID, work)
			if res.StatusCode == 409 { // duplicate
				continue
			} else if err != nil {
				return err
			}

			pub.ORCIDWork = append(pub.ORCIDWork, models.PublicationORCIDWork{
				ORCID:   args.ORCID,
				PutCode: putCode,
			})

			if _, err := a.PublicationService.UpdatePublication(pub); err != nil {
				return err
			}
		}

		activity.RecordHeartbeat(ctx, models.Progress{Numerator: numDone, Denominator: hits.Total})

		if !hits.NextPage {
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

	return w
}
