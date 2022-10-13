package commands

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/ulid"
)

func init() {
	publicationGetCmd.Flags().StringP("format", "f", "jsonl", "export format")
	publicationAddCmd.Flags().StringP("format", "f", "jsonl", "import format")
	publicationCmd.AddCommand(publicationGetCmd)
	publicationCmd.AddCommand(publicationAllCmd)
	publicationCmd.AddCommand(publicationAddCmd)
	publicationCmd.AddCommand(publicationImportCmd)
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
			if err := e.Repository.SavePublication(&p); err != nil {
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
				log.Printf("Validation failed for publication [id: %s] at line %d : %v", p.ID, lineNo, err)
				continue
			}
			if err := e.Repository.ImportPublication(&p); err != nil {
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
