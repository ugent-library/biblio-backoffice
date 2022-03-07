package orcid

import (
	"fmt"
	"strings"
	"time"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-orcid/orcid"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/text/language"
)

type AddpublicationsWorkflowArgs struct {
	UserID     string
	ORCID      string
	ORCIDToken string
	SearchArgs models.SearchArgs
}

type AddpublicationsArgs struct {
	ORCID      string
	ORCIDToken string
	Hits       models.PublicationHits
}

func AddPublicationsWorkflow(publicationSearchService backends.PublicationSearchService) func(workflow.Context, AddpublicationsWorkflowArgs) error {
	return func(ctx workflow.Context, args AddpublicationsWorkflowArgs) error {
		exportRetrypolicy := &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1.0,
			MaximumInterval:    time.Second * 10, // 10 * InitialInterval
			MaximumAttempts:    3,                // Do it for a minute
			//NonRetryableErrorTypes: []string, // empty
		}

		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
			RetryPolicy:         exportRetrypolicy,
		})

		taskState := models.TaskState{
			Message: "Adding publications to your ORCID works",
			Status:  models.Waiting,
		}

		logger := workflow.GetLogger(ctx)
		// setup query handler for query type "state"
		err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (models.TaskState, error) {
			return taskState, nil
		})
		if err != nil {
			logger.Info("SetQueryHandler failed: " + err.Error())
			return err
		}

		taskState.Status = models.Running

		searchArgs := args.SearchArgs

		for {
			hits, _ := publicationSearchService.UserPublications(args.UserID, &searchArgs)

			taskState.Denominator = hits.Total

			logger.Info("execute activity")

			err = workflow.ExecuteActivity(ctx, "AddPublicationsToORCID", AddpublicationsArgs{
				ORCID:      args.ORCID,
				ORCIDToken: args.ORCIDToken,
				Hits:       *hits,
			}).Get(ctx, nil)
			if err != nil {
				taskState.Message = "Adding publications to your ORCID works failed"
				taskState.Status = models.Failed
				logger.Error("AddPublicationsToORCID failed.", "Error", err)
				return err
			}

			taskState.Numerator += len(hits.Hits)

			if !hits.NextPage {
				taskState.Message = fmt.Sprintf("Added %d publications to your ORCID works", hits.Total)
				taskState.Status = models.Done
				break
			}
			searchArgs.Page = searchArgs.Page + 1
		}

		return nil
	}
}

func AddPublications(publicationService backends.PublicationService, orcidSandbox bool) func(AddpublicationsArgs) error {
	return func(args AddpublicationsArgs) error {
		// logger := workflow.GetLogger(ctx)
		// logger.Info("sending pubs to orcid")

		orcidClient := orcid.NewMemberClient(orcid.Config{
			Token:   args.ORCIDToken,
			Sandbox: orcidSandbox,
		})

		for _, pub := range args.Hits.Hits {
			var done bool
			for _, ow := range pub.ORCIDWork {
				if ow.ORCID == args.ORCID { // already sent to orcid
					done = true
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

			if _, err := publicationService.UpdatePublication(pub); err != nil {
				return err
			}
		}

		return nil
	}
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
