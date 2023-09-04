package commands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/oai-service/api/v1"

	"github.com/ugent-library/biblio-backoffice/internal/backends/oaidc"
	"github.com/ugent-library/biblio-backoffice/internal/models"
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
		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "biblio:fulltext",
			SetName: "Biblio records with a fulltext file",
		})
		if err != nil {
			logger.Fatal(err)
		}
		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "biblio:open_access",
			SetName: "Biblio records with an open access fulltext file",
		})
		if err != nil {
			logger.Fatal(err)
		}

		// add all publications
		repo := Services().Repository
		repo.EachPublication(func(p *models.Publication) bool {
			oaiID := "oai:archive.ugent.be:" + p.ID

			if p.Status == "deleted" && p.HasBeenPublic {
				err = client.DeleteRecord(context.TODO(), &api.DeleteRecordRequest{
					Identifier: oaiID,
				})
				if err != nil {
					logger.Fatal(err)
				}
				return true
			}

			if p.Status != "public" {
				return true
			}

			metadata, err := oaidc.EncodePublication(p)
			if err != nil {
				logger.Fatal(err)
			}

			err = client.AddRecord(context.TODO(), &api.AddRecordRequest{
				Identifier:     oaiID,
				MetadataPrefix: "oai_dc",
				Content:        string(metadata),
			})
			if err != nil {
				logger.Fatal(err)
			}

			setSpecs := []string{"biblio"}

			if p.Type == "journal_article" {
				setSpecs = append(setSpecs, "biblio:journal_article")
			}
			if p.Type == "book" {
				setSpecs = append(setSpecs, "biblio:book")
			}
			for _, f := range p.File {
				if f.Relation == "main_file" {
					setSpecs = append(setSpecs, "biblio:fulltext")
					break
				}
			}
			for _, f := range p.File {
				if f.Relation == "main_file" && f.AccessLevel == "info:eu-repo/semantics/openAccess" {
					setSpecs = append(setSpecs, "biblio:open_access")
					break
				}
			}

			err = client.AddItem(context.TODO(), &api.AddItemRequest{
				Identifier: oaiID,
				SetSpecs:   setSpecs,
			})
			if err != nil {
				logger.Fatal(err)
			}

			return true
		})
	},
}
