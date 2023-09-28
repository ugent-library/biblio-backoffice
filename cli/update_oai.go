package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/oai-service/api/v1"

	"github.com/ugent-library/biblio-backoffice/backends/mods36"
	"github.com/ugent-library/biblio-backoffice/backends/oaidc"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		oaiEncoder := oaidc.New(config.Frontend.URL)
		modsEncoder := mods36.New(config.Frontend.URL)

		client, err := api.NewClient(config.OAI.APIURL, &securitySource{config.OAI.APIKey})
		if err != nil {
			return err
		}

		err = client.AddMetadataFormat(context.TODO(), &api.AddMetadataFormatRequest{
			MetadataPrefix:    "oai_dc",
			MetadataNamespace: "http://www.openarchives.org/OAI/2.0/oai_dc/",
			Schema:            "http://www.openarchives.org/OAI/2.0/oai_dc.xsd",
		})
		if err != nil {
			return err
		}
		err = client.AddMetadataFormat(context.TODO(), &api.AddMetadataFormatRequest{
			MetadataPrefix:    "mods_36",
			MetadataNamespace: "http://www.loc.gov/mods/v3",
			Schema:            "http://www.loc.gov/standards/mods/v3/mods-3-6.xsd",
		})
		if err != nil {
			return err
		}

		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "biblio",
			SetName: "All Biblio records",
		})
		if err != nil {
			return err
		}
		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "biblio:fulltext",
			SetName: "Biblio records with a fulltext file",
		})
		if err != nil {
			return err
		}
		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "biblio:open_access",
			SetName: "Biblio records with an open access fulltext file",
		})
		if err != nil {
			logger.Fatal(err)
		}
		for _, t := range vocabularies.Map["publication_types"] {
			err = client.AddSet(context.TODO(), &api.AddSetRequest{
				SetSpec: "biblio:" + t,
				SetName: "Biblio " + t + " records",
			})
			if err != nil {
				return err
			}
		}

		// add all publications
		repo := newServices().Repo
		repo.EachPublication(func(p *models.Publication) bool {
			oaiID := "oai:archive.ugent.be:" + p.ID

			if p.Status == "deleted" && p.HasBeenPublic {
				err = client.DeleteRecord(context.TODO(), &api.DeleteRecordRequest{
					Identifier: oaiID,
				})
				if err != nil {
					// TODO
					logger.Fatal(err)
				}
				return true
			}

			if p.Status != "public" {
				return true
			}

			metadata, err := oaiEncoder.EncodePublication(p)
			if err != nil {
				// TODO
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

			metadata, err = modsEncoder.EncodePublication(p)
			if err != nil {
				// TODO
				logger.Fatal(err)
			}

			err = client.AddRecord(context.TODO(), &api.AddRecordRequest{
				Identifier:     oaiID,
				MetadataPrefix: "mods_36",
				Content:        string(metadata),
			})
			if err != nil {
				// TODO
				logger.Fatal(err)
			}

			setSpecs := []string{"biblio", "biblio:" + p.Type}

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
				// TODO
				logger.Fatal(err)
			}

			return true
		})

		return nil
	},
}
