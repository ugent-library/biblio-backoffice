package cli

import (
	"context"
	"regexp"
	"slices"

	"github.com/spf13/cobra"
	"github.com/ugent-library/oai-service/api/v1"

	"github.com/ugent-library/biblio-backoffice/backends/mods36"
	"github.com/ugent-library/biblio-backoffice/backends/oaidc"
	"github.com/ugent-library/biblio-backoffice/models"
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
		services := newServices()

		reFP := regexp.MustCompile(`^FP[0-9]+$`)

		oaiEncoder := oaidc.New(services.Repo, config.Frontend.URL)
		modsEncoder := mods36.New(services.Repo, config.Frontend.URL)

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
			SetSpec: "fulltext",
			SetName: "Biblio records with a fulltext file",
		})
		if err != nil {
			return err
		}

		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "open_access",
			SetName: "Biblio records with an open access fulltext file",
		})
		if err != nil {
			return err
		}

		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "ec_fundedresources",
			SetName: "OpenAire 2.0",
		})
		if err != nil {
			return err
		}

		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "openaire",
			SetName: "OpenAire 3.0",
		})
		if err != nil {
			return err
		}

		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "driver",
			SetName: "Driver",
		})
		if err != nil {
			return err
		}

		err = client.AddSet(context.TODO(), &api.AddSetRequest{
			SetSpec: "iminds",
			SetName: "All iMinds publications",
		})
		if err != nil {
			return err
		}

		repo := newServices().Repo

		// add all publications
		repo.EachPublication(func(p *models.Publication) bool {
			oaiID := "oai:archive.ugent.be:" + p.ID

			if p.HasBeenPublic && p.Status != "public" {
				for _, metadataPrefix := range []string{"oai_dc", "mods_36"} {
					err = client.DeleteRecord(context.TODO(), &api.DeleteRecordRequest{
						Identifier:     oaiID,
						MetadataPrefix: metadataPrefix,
					})
					if err != nil {
						logger.Error("cannot delete oai publication", "identifier", oaiID, "metadataPrefix", metadataPrefix)
					}
				}
				return true
			}

			if p.Status != "public" {
				return true
			}

			metadata, err := oaiEncoder.EncodePublication(p)
			if err != nil {
				logger.Error("cannot encode oai publication", "identifier", oaiID, "metadataPrefix", "oai_dc")
				return true
			}

			err = client.AddRecord(context.TODO(), &api.AddRecordRequest{
				Identifier:     oaiID,
				MetadataPrefix: "oai_dc",
				Content:        string(metadata),
			})
			if err != nil {
				logger.Error("cannot add oai publication", "identifier", oaiID, "metadataPrefix", "oai_dc")
				return true
			}

			metadata, err = modsEncoder.EncodePublication(p)
			if err != nil {
				logger.Error("cannot encode oai publication", "identifier", oaiID, "metadataPrefix", "mods_36")
				return true
			}

			err = client.AddRecord(context.TODO(), &api.AddRecordRequest{
				Identifier:     oaiID,
				MetadataPrefix: "mods_36",
				Content:        string(metadata),
			})
			if err != nil {
				logger.Error("cannot add oai publication", "identifier", oaiID, "metadataPrefix", "mods_36")
				return true
			}

			setSpecs := []string{}

			for _, f := range p.File {
				if f.Relation == "main_file" {
					setSpecs = append(setSpecs, "fulltext")
					break
				}
			}
			for _, f := range p.File {
				if f.Relation == "main_file" && f.AccessLevel == "info:eu-repo/semantics/openAccess" {
					setSpecs = append(setSpecs, "open_access", "driver")
					break
				}
			}

			for _, rp := range p.RelatedProjects {
				if rp.Project.EUProject != nil && (rp.Project.EUProject.FrameworkProgramme == "H2020" || reFP.MatchString(rp.Project.EUProject.FrameworkProgramme)) {
					setSpecs = append(setSpecs, "ec_fundedresources")
					break
				}
			}

			if slices.Contains(setSpecs, "open_access") || slices.Contains(setSpecs, "ec_fundedresources") {
				setSpecs = append(setSpecs, "openaire")
			}

			for _, relOrg := range p.RelatedOrganizations {
				if relOrg.OrganizationID == "IBBT" {
					setSpecs = append(setSpecs, "iminds")
					break
				}
			}

			err = client.AddItem(context.TODO(), &api.AddItemRequest{
				Identifier: oaiID,
				SetSpecs:   setSpecs,
			})
			if err != nil {
				logger.Error("cannot add oai dataset item", "identifier", oaiID)
				return true
			}

			return true
		})

		repo.EachDataset(func(d *models.Dataset) bool {
			oaiID := "oai:archive.ugent.be:" + d.ID

			// any record that has been part of the oai set but is now not visible anymore
			// because it is returned or deleted is marked as deleted in oai
			if d.HasBeenPublic && d.Status != "public" {
				for _, metadataPrefix := range []string{"oai_dc", "mods_36"} {
					err = client.DeleteRecord(context.TODO(), &api.DeleteRecordRequest{
						Identifier:     oaiID,
						MetadataPrefix: metadataPrefix,
					})
					if err != nil {
						logger.Error("cannot delete oai dataset", "identifier", oaiID, "metadataPrefix", metadataPrefix)
					}
				}
				return true
			}

			if d.Status != "public" {
				return true
			}

			metadata, err := oaiEncoder.EncodeDataset(d)
			if err != nil {
				logger.Error("cannot encode oai dataset", "identifier", oaiID, "metadataPrefix", "oai_dc")
				return true
			}

			err = client.AddRecord(context.TODO(), &api.AddRecordRequest{
				Identifier:     oaiID,
				MetadataPrefix: "oai_dc",
				Content:        string(metadata),
			})
			if err != nil {
				logger.Error("cannot add oai dataset", "identifier", oaiID, "metadataPrefix", "oai_dc")
				return true
			}

			metadata, err = modsEncoder.EncodeDataset(d)
			if err != nil {
				logger.Error("cannot encode oai dataset", "identifier", oaiID, "metadataPrefix", "mods_36")
				return true
			}

			err = client.AddRecord(context.TODO(), &api.AddRecordRequest{
				Identifier:     oaiID,
				MetadataPrefix: "mods_36",
				Content:        string(metadata),
			})
			if err != nil {
				logger.Error("cannot add oai dataset", "identifier", oaiID, "metadataPrefix", "mods_36")
				return true
			}

			err = client.AddItem(context.TODO(), &api.AddItemRequest{
				Identifier: oaiID,
				SetSpecs:   []string{},
			})
			if err != nil {
				logger.Error("cannot add oai dataset item", "identifier", oaiID)
				return true
			}

			return true
		})

		return nil
	},
}
