package commands

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/internal/backends/oaidc"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

func init() {
	rootCmd.AddCommand(updateOai)
}

var updateOai = &cobra.Command{
	Use:   "update-oai",
	Short: "Update OAI provider",
	Run: func(cmd *cobra.Command, args []string) {
		logger := newLogger()
		// setup nats connection
		nc, err := nats.Connect("nats://localhost:4222")
		if err != nil {
			logger.Fatal(err)
		}
		// add all publications
		repo := Services().Repository
		repo.EachPublication(func(p *models.Publication) bool {
			md, err := oaidc.EncodePublication(p)
			if err != nil {
				logger.Fatal(err)
			}
			req := struct {
				Identifier     string   `json:"identifier"`
				MetadataPrefix string   `json:"metadata_prefix"`
				Metadata       string   `json:"metadata"`
				SetSpecs       []string `json:"set_specs"`
			}{
				Identifier:     p.ID,
				MetadataPrefix: "oai_dc",
				Metadata:       string(md),
			}
			data, err := json.Marshal(req)
			res, err := nc.Request("oai.AddRecord", data, 1*time.Second)
			if err != nil {
				logger.Fatal(err)
			}
			logger.Info(res)
			return true
		})
	},
}
