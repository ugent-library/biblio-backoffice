package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/ulid"
)

func init() {
	publicationGetCmd.Flags().StringP("format", "f", "jsonl", "export format")
	publicationAddCmd.Flags().StringP("format", "f", "jsonl", "import format")
	publicationImportCmd.Flags().StringP("format", "f", "jsonl", "import format")
	oldPublicationImportCmd.Flags().StringP("format", "f", "jsonl", "import format")
	publicationCmd.AddCommand(publicationGetCmd)
	publicationCmd.AddCommand(publicationAllCmd)
	publicationCmd.AddCommand(publicationAddCmd)
	publicationCmd.AddCommand(publicationImportCmd)
	publicationCmd.AddCommand(oldPublicationImportCmd)
	publicationCmd.AddCommand(updatePublicationEmbargoes)
	publicationCmd.AddCommand(createPublicationHandles)
	publicationCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(publicationCmd)
}

var publicationCmd = &cobra.Command{
	Use:   "publication [command]",
	Short: "Publication commands",
}

var publicationGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get publication by id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()

		fmt, _ := cmd.Flags().GetString("format")

		if fmt != "jsonl" {
			enc, ok := e.PublicationEncoders[fmt]
			if !ok {
				log.Fatalf("Unknown format %s", fmt)
			}
			d, err := e.Repository.GetPublication(args[0])
			if err != nil {
				log.Fatal(err)
			}
			b, err := enc(d)
			if err != nil {
				log.Fatal(err)

			}
			os.Stdout.Write(b)
			return
		}

		enc := json.NewEncoder(os.Stdout)
		d, err := e.Repository.GetPublication(args[0])
		if err != nil {
			log.Fatal(err)
		}
		enc.Encode(d)
	},
}

var publicationAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Get all publications",
	Run: func(cmd *cobra.Command, args []string) {
		s := newRepository()
		e := json.NewEncoder(os.Stdout)
		s.EachPublication(func(d *models.Publication) bool {
			e.Encode(d)
			return true
		})
	},
}

var publicationAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add publications",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()

		var indexWG sync.WaitGroup

		// indexing channel
		indexC := make(chan *models.Publication)

		// start bulk indexer
		indexWG.Add(1)
		go func() {
			defer indexWG.Done()
			e.PublicationSearchService.IndexMultiple(indexC)
		}()

		fmt, _ := cmd.Flags().GetString("format")
		decFactory, ok := e.PublicationDecoders[fmt]
		if !ok {
			log.Fatalf("Unknown format %s", fmt)
		}
		dec := decFactory(os.Stdin)

		lineNo := 0
		for {
			lineNo += 1
			p := models.Publication{
				ID:             ulid.MustGenerate(),
				Status:         "private",
				Classification: "U",
			}
			if err := dec.Decode(&p); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				log.Fatalf("Unable to decode publication at line %d : %v", lineNo, err)
			}
			if err := p.Validate(); err != nil {
				log.Printf("Validation failed for publication [id: %s] at line %d : %v", p.ID, lineNo, err)
				continue
			}
			if err := e.Repository.SavePublication(&p, nil); err != nil {
				log.Fatalf("Unable to store publication from line %d : %v", lineNo, err)
			}

			indexC <- &p
		}

		// close indexing channel when all recs are stored
		close(indexC)
		// wait for indexing to finish
		indexWG.Wait()
	},
}

var publicationImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import publications",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()

		var indexWG sync.WaitGroup

		// indexing channel
		indexC := make(chan *models.Publication)

		// start bulk indexer
		indexWG.Add(1)
		go func() {
			defer indexWG.Done()
			e.PublicationSearchService.IndexMultiple(indexC)
		}()

		fmt, _ := cmd.Flags().GetString("format")
		decFactory, ok := e.PublicationDecoders[fmt]
		if !ok {
			log.Fatalf("Unknown format %s", fmt)
		}
		dec := decFactory(os.Stdin)

		lineNo := 0
		for {
			lineNo += 1
			p := models.Publication{}
			if err := dec.Decode(&p); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				log.Fatalf("Unable to decode publication at line %d : %v", lineNo, err)
			}
			if err := p.Validate(); err != nil {
				log.Printf(
					"Validation failed for publication[snapshot_id: %s, id: %s] at line %d : %v",
					p.SnapshotID,
					p.ID,
					lineNo,
					err,
				)
				continue
			}
			if err := e.Repository.ImportCurrentPublication(&p); err != nil {
				log.Printf(
					"Unable to store publication[snapshot_id: %s, id: %s] from line %d : %v",
					p.SnapshotID,
					p.ID,
					lineNo,
					err,
				)
				continue
			}
			log.Printf(
				"Added publication[snapshot_id: %s, id: %s]",
				p.SnapshotID,
				p.ID,
			)

			indexC <- &p
		}

		// close indexing channel when all recs are stored
		close(indexC)
		// wait for indexing to finish
		indexWG.Wait()
	},
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Make publications consistent, clean up data anomalies",
	Run: func(cmd *cobra.Command, args []string) {
		s := newRepository()
		e := Services()

		var indexWG sync.WaitGroup

		// indexing channel
		indexC := make(chan *models.Publication)

		// start bulk indexer
		indexWG.Add(1)
		go func() {
			defer indexWG.Done()
			e.PublicationSearchService.IndexMultiple(indexC)
		}()

		s.EachPublication(func(p *models.Publication) bool {
			// Guard
			fixed := false

			// Make sure to set p.User = nil!!
			p.User = nil

			// Add the department "tree" property if it is missing.
			for _, dep := range p.Department {
				if dep.Tree == nil {
					depID := dep.ID
					org, orgErr := e.OrganizationService.GetOrganization(depID)
					if orgErr == nil {
						p.RemoveDepartment(depID)
						p.AddDepartmentByOrg(org)
						fixed = true
					}
				}
			}

			if fixed {
				if err := p.Validate(); err != nil {
					log.Printf(
						"Validation failed for publication[snapshot_id: %s, id: %s] : %v",
						p.SnapshotID,
						p.ID,
						err,
					)
					return false
				}

				err := s.UpdatePublication(p.SnapshotID, p, nil)

				var conflict *snapstore.Conflict
				if errors.As(err, &conflict) {
					log.Printf(
						"Conflict detected for publication[snapshot_id: %s, id: %s] : %v",
						p.SnapshotID,
						p.ID,
						err,
					)
					return false
				}

				log.Printf(
					"Fixed publication[snapshot_id: %s, id: %s]",
					p.SnapshotID,
					p.ID,
				)

				indexC <- p
			}

			return true
		})

		// close indexing channel when all recs are stored
		close(indexC)
		// wait for indexing to finish
		indexWG.Wait()
	},
}

var oldPublicationImportCmd = &cobra.Command{
	Use:   "import-version",
	Short: "Import old publications",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()

		fmt, _ := cmd.Flags().GetString("format")
		decFactory, ok := e.PublicationDecoders[fmt]
		if !ok {
			log.Fatalf("Unknown format %s", fmt)
		}
		dec := decFactory(os.Stdin)

		lineNo := 0
		for {
			lineNo += 1
			p := models.Publication{}
			if err := dec.Decode(&p); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				log.Fatalf("Unable to decode publication at line %d : %v", lineNo, err)
			}
			if err := p.Validate(); err != nil {
				log.Printf(
					"Validation failed for publication[snapshot_id: %s, id: %s] at line %d : %v",
					p.SnapshotID,
					p.ID,
					lineNo,
					err,
				)
				continue
			}
			if err := e.Repository.ImportOldPublication(&p); err != nil {
				log.Printf(
					"Unable to store old publication[snapshot_id: %s, id: %s] from line %d : %v",
					p.SnapshotID,
					p.ID,
					lineNo,
					err,
				)
				continue
			}
			log.Printf(
				"Added old publication[snapshot_id: %s, id: %s]",
				p.SnapshotID,
				p.ID,
			)
		}
	},
}

var updatePublicationEmbargoes = &cobra.Command{
	Use:   "update-embargoes",
	Short: "Update publication embargoes",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()
		logger := newLogger()

		e.Repository.AddPublicationListener(func(p *models.Publication) {
			if err := e.PublicationSearchService.Index(p); err != nil {
				logger.Fatalf("error indexing publication %s: %v", p.ID, err)
			}
		})

		var count int = 0
		updateEmbargoErr := e.Repository.Transaction(
			context.Background(),
			func(repo backends.Repository) error {

				/*
					select live publications that have files with embargoed access
				*/
				var embargoAccessLevel string = "info:eu-repo/semantics/embargoedAccess"
				currentDateStr := time.Now().Format("2006-01-02")
				var sqlPublicationWithEmbargo string = `
				SELECT * FROM publications WHERE date_until IS NULL AND
				data->'file' IS NOT NULL AND
				EXISTS(
					SELECT 1 FROM jsonb_array_elements(data->'file') AS f
					WHERE f->>'access_level' = $1 AND
					f->>'embargo_date' <= $2
				)
				`

				publications := make([]*models.Publication, 0)
				sErr := repo.SelectPublications(
					sqlPublicationWithEmbargo,
					[]any{
						embargoAccessLevel,
						currentDateStr},
					func(publication *models.Publication) bool {
						publications = append(publications, publication)
						return true
					},
				)

				if sErr != nil {
					return sErr
				}

				for _, publication := range publications {
					/*
						clear outdated embargoes
					*/
					for _, file := range publication.File {
						if file.AccessLevel != embargoAccessLevel {
							continue
						}
						// TODO: what with empty embargo_date?
						if file.EmbargoDate == "" {
							continue
						}
						if file.EmbargoDate > currentDateStr {
							continue
						}
						file.ClearEmbargo()
					}

					publication.User = nil

					if e := repo.SavePublication(publication, nil); e != nil {
						return e
					}
					count++
				}

				return nil
			},
		)

		if updateEmbargoErr != nil {
			logger.Fatal(updateEmbargoErr)
		}

		logger.Infof("updated %d embargoes", count)
	},
}

var createPublicationHandles = &cobra.Command{
	Use:   "create-handles",
	Short: "Create publication handles",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()
		logger := newLogger()
		handleService := e.HandleService

		if handleService == nil {
			logger.Fatal("handle server updates are not enabled")
		}

		e.Repository.AddPublicationListener(func(p *models.Publication) {
			if err := e.PublicationSearchService.Index(p); err != nil {
				logger.Fatalf("error indexing publication %s: %v", p.ID, err)
			}
		})

		var count int = 0
		createHandlesErr := e.Repository.Transaction(
			context.Background(),
			func(repo backends.Repository) error {

				publications := make([]*models.Publication, 0)
				sql := `
				SELECT * FROM publications WHERE date_until IS NULL AND
				data->>'status' = 'public' AND
				NOT data ? 'handle'
				`

				selectErr := repo.SelectPublications(
					sql,
					[]any{},
					func(p *models.Publication) bool {
						publications = append(publications, p)
						return true
					},
				)

				if selectErr != nil {
					return selectErr
				}

				for _, p := range publications {
					h, hErr := handleService.UpsertHandleByPublication(p)
					if hErr != nil {
						return fmt.Errorf(
							"error adding handle for publication %s: %s",
							p.ID,
							hErr,
						)
					} else if !h.IsSuccess() {
						return fmt.Errorf(
							"error adding handle for publication %s: %s",
							p.ID,
							h.Message,
						)
					} else {
						logger.Infof(
							"added handle url %s to publication %s",
							h.GetFullHandleURL(),
							p.ID,
						)
						p.Handle = h.GetFullHandleURL()
						p.User = nil

						if e := repo.SavePublication(p, nil); e != nil {
							return e
						}
						count++
					}
				}

				return nil
			},
		)

		if createHandlesErr != nil {
			logger.Fatal(createHandlesErr)
		}

		logger.Infof("created %d handles", count)
	},
}
