package commands

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/snapstore"
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
	publicationCmd.AddCommand(cleanupCmd)
	publicationCmd.AddCommand(transferCmd)
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
				ID:             ulid.Make().String(),
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

			// Trim keywords, remove empty keywords
			var cleanKeywords []string
			for _, kw := range p.Keyword {
				cleanKw := strings.TrimSpace(kw)
				if cleanKw != kw || cleanKw == "" {
					fixed = true
				}
				if cleanKw != "" {
					cleanKeywords = append(cleanKeywords, cleanKw)
				}
			}
			p.Keyword = cleanKeywords

			// Save record if changed
			if fixed {
				p.User = nil

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

var transferCmd = &cobra.Command{
	Use:   "transfer UID UID [PUBID]",
	Short: "Transfer publications between people",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()
		s := newRepository()

		s.AddPublicationListener(func(p *models.Publication) {
			if p.DateUntil == nil {
				if err := e.PublicationSearchService.Index(p); err != nil {
					log.Fatalf("error indexing publication %s: %v", p.ID, err)
				}
			}
		})

		source := args[0]
		dest := args[1]

		p, err := e.PersonService.GetPerson(dest)
		if err != nil {
			log.Printf("Fatal: could not retrieve person %s: %s", dest, err)
		}

		c := &models.Contributor{}
		c.ID = p.ID
		c.FirstName = p.FirstName
		c.LastName = p.LastName
		c.FullName = p.FullName
		c.UGentID = p.UGentID
		c.ORCID = p.ORCID
		for _, pd := range p.Department {
			newDep := models.ContributorDepartment{ID: pd.ID}
			org, orgErr := e.OrganizationService.GetOrganization(pd.ID)
			if orgErr == nil {
				newDep.Name = org.Name
			}
			c.Department = append(c.Department, newDep)
		}

		callback := func(p *models.Publication) bool {
			fixed := false

			if p.User != nil {
				if p.User.ID == source {
					p.User = &models.PublicationUser{
						ID:   c.ID,
						Name: c.FullName,
					}

					log.Printf("p: %s: s: %s ::: user: %s -> %s", p.ID, p.SnapshotID, source, c.ID)
					fixed = true
				}
			}

			if p.Creator != nil {
				if p.Creator.ID == source {
					p.Creator = &models.PublicationUser{
						ID:   c.ID,
						Name: c.FullName,
					}

					if len(c.Department) > 0 {
						org, orgErr := e.OrganizationService.GetOrganization(c.Department[0].ID)
						if orgErr != nil {
							log.Printf("p: %s: s: %s ::: creator: could not fetch department for %s: %s", p.ID, p.SnapshotID, c.ID, orgErr)
						} else {
							p.AddDepartmentByOrg(org)
						}
					}

					log.Printf("p: %s: s: %s ::: creator: %s -> %s", p.ID, p.SnapshotID, source, c.ID)
					fixed = true
				}
			}

			for k, a := range p.Author {
				if a.ID == source {
					p.SetContributor("author", k, c)
					log.Printf("p: %s: s: %s ::: author: %s -> %s", p.ID, p.SnapshotID, a.ID, c.ID)
					fixed = true
				}
			}

			for k, e := range p.Editor {
				if e.ID == source {
					p.SetContributor("editor", k, c)
					log.Printf("p: %s: s: %s ::: editor: %s -> %s", p.ID, p.SnapshotID, e.ID, c.ID)
					fixed = true
				}
			}

			for k, s := range p.Supervisor {
				if s.ID == source {
					p.SetContributor("supervisor", k, c)
					log.Printf("p: %s: s: %s ::: supervisor: %s -> %s", p.ID, p.SnapshotID, s.ID, c.ID)
					fixed = true
				}
			}

			if fixed {
				errUpdate := s.UpdatePublicationInPlace(p)
				if errUpdate != nil {
					log.Printf("p: %s: s: %s ::: Could not update snapshot: %s", p.ID, p.SnapshotID, errUpdate)
				}
			}

			return true
		}

		if len(args) > 2 {
			pubID := args[2]
			s.PublicationHistory(pubID, callback)
		} else {
			s.EachPublicationSnapshot(callback)
		}
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
