package commands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backoffice/internal/backends/oaidc"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/oai-service/api/v1"
)

func init() {
	rootCmd.AddCommand(updateOai)
}

type securitySource struct {
	apiKey string
}

func (s *securitySource) ApiKey(ctx context.Context, operationName string) (api.ApiKey, error) {
	return api.ApiKey{APIKey: s.apiKey}, nil
}

var updateOai = &cobra.Command{
	Use:   "update-oai",
	Short: "Update OAI provider",
	Run: func(cmd *cobra.Command, args []string) {
		logger := newLogger()

		client, err := api.NewClient(viper.GetString("oai-api-url"), &securitySource{viper.GetString("oai-api-key")})
		if err != nil {
			logger.Fatal(err)
		}

		// TODO
		err = client.AddMetadataFormat(context.TODO(), &api.AddMetadataFormatRequest{
			MetadataPrefix:    "oai_dc",
			MetadataNamespace: "http://www.openarchives.org/OAI/2.0/oai_dc/",
			Schema:            "http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
		})
		if err != nil {
			logger.Fatal(err)
		}
		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "biblio",
			SetName: "All Biblio records",
		})
		if err != nil {
			logger.Fatal(err)
		}
		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "biblio:journal_article",
			SetName: "Biblio journal articles",
		})
		if err != nil {
			logger.Fatal(err)
		}
		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "biblio:book",
			SetName: "Biblio books",
		})
		if err != nil {
			logger.Fatal(err)
		}

		// add all publications
		n := 0
		repo := Services().Repository
		repo.EachPublication(func(p *models.Publication) bool {
			if p.Status != "public" {
				return true
			}

			metadata, err := oaidc.EncodePublication(p)
			if err != nil {
				logger.Fatal(err)
			}

			err = client.AddRecordMetadata(context.TODO(), &api.AddRecordMetadataRequest{
				Identifier:     p.ID,
				MetadataPrefix: "oai_dc",
				Content:        string(metadata),
			})
			if err != nil {
				logger.Fatal(err)
			}

			if p.Type == "journal_article" {
				err = client.AddRecordSets(context.TODO(), &api.AddRecordSetsRequest{
					Identifier: p.ID,
					SetSpecs:   []string{"biblio:journal_article"},
				})
				if err != nil {
					logger.Fatal(err)
				}
			}
			if p.Type == "book" {
				err = client.AddRecordSets(context.TODO(), &api.AddRecordSetsRequest{
					Identifier: p.ID,
					SetSpecs:   []string{"biblio:book"},
				})
				if err != nil {
					logger.Fatal(err)
				}
			}

			n++

			return n < 5000
		})
	},
}
