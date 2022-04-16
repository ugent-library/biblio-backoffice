package commands

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/models"
)

func init() {
	publicationGetCmd.Flags().StringP("format", "f", "", "export format")

	publicationCmd.AddCommand(publicationGetCmd)
	publicationCmd.AddCommand(publicationAllCmd)
	publicationCmd.AddCommand(publicationAddCmd)
	rootCmd.AddCommand(publicationCmd)
}

var publicationCmd = &cobra.Command{
	Use:   "publication [command]",
	Short: "Publication commands",
}

var publicationGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get publication by id",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		e := Engine()

		if format := viper.GetString("format"); format != "" {
			enc, ok := e.PublicationEncoders[format]
			if !ok {
				log.Fatalf("Unknown format %s", format)
			}
			d, err := e.StorageService.GetPublication(args[0])
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
		d, err := e.StorageService.GetPublication(args[0])
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
		s := newStorageService()
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
		e := Engine()

		var indexWG sync.WaitGroup

		// indexing channel
		indexC := make(chan *models.Publication)

		// start bulk indexer
		go func() {
			indexWG.Add(1)
			defer indexWG.Done()
			e.PublicationSearchService.IndexPublications(indexC)
		}()

		dec := json.NewDecoder(os.Stdin)
		for {
			var p models.Publication
			if err := dec.Decode(&p); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			if err := p.Validate(); err != nil {
				log.Fatal(err)
			}

			savedP, err := e.StorageService.SavePublication(&p)
			if err != nil {
				log.Fatal(err)
			}

			indexC <- savedP
		}

		// close indexing channel when all recs are stored
		close(indexC)
		// wait for indexing to finish
		indexWG.Wait()
	},
}
