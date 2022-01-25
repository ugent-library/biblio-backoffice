package orcidworker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/rabbitmq/amqp091-go"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-orcid/orcid"
	"golang.org/x/text/language"
)

type service struct {
	engine  *engine.Engine
	mqChan  *amqp091.Channel
	msgChan <-chan amqp091.Delivery
}

type addToORCID struct {
	UserID     string             `json:"user_id"`
	SearchArgs *engine.SearchArgs `json:"search_args"`
}

func New(e *engine.Engine) (*service, error) {
	mqCh, err := e.MQ.Channel()
	if err != nil {
		log.Fatal(err)
	}

	_, err = mqCh.QueueDeclare(
		"tasks_orcid", // queue name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return nil, err
	}

	err = mqCh.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}

	msgChan, err := mqCh.Consume(
		"tasks_orcid", // queue
		"",            // consumer
		false,         // auto ack
		false,         // exclusive
		false,         // no local
		false,         // no wait
		nil,           // args
	)
	if err != nil {
		return nil, err
	}

	return &service{e, mqCh, msgChan}, nil
}

func (s *service) Name() string {
	return "orcid-worker"
}

func (s *service) Start() error {
	go func() {
		for d := range s.msgChan {
			task := addToORCID{}
			if err := json.Unmarshal(d.Body, &task); err != nil {
				log.Println(err)
			}
			s.handleTask(task)
			d.Ack(false)
		}
	}()

	return nil
}

func (s *service) Stop(ctx context.Context) error {
	return s.mqChan.Close()
}

func (s *service) handleTask(task addToORCID) error {
	user, _ := s.engine.GetUser(task.UserID)
	args := task.SearchArgs
	client := orcid.NewMemberClient(orcid.Config{
		Token:   user.ORCIDToken,
		Sandbox: s.engine.ORCIDSandbox,
	})

	for {
		hits, _ := s.engine.UserPublications(user.ID, args)
		for _, pub := range hits.Hits {
			for _, ow := range pub.ORCIDWork {
				if ow.ORCID == user.ORCID {
					continue
				}
			}

			work := publicationToORCID(pub)
			putCode, res, err := client.AddWork(user.ORCID, work)
			if err != nil {
				log.Println(err)
				body, _ := ioutil.ReadAll(res.Body)
				log.Print(string(body))
				continue
			}

			pub.ORCIDWork = append(pub.ORCIDWork, models.PublicationORCIDWork{
				ORCID:   user.ORCID,
				PutCode: putCode,
			})

			if _, err := s.engine.UpdatePublication(pub); err != nil {
				log.Println(err)
			}

			log.Printf("processed task orcid add for pub %s and user %s", pub.ID, user.ID)
		}
		if !hits.NextPage {
			break
		}
		args.Page = args.Page + 1
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

	log.Printf("%+v", w)

	return w
}
