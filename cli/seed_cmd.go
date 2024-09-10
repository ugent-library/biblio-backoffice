package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/backends/authority"
	"github.com/ugent-library/biblio-backoffice/models"
)

func init() {
	rootCmd.AddCommand(seedCmd)

	seedOrganizationsCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.AddCommand(seedOrganizationsCmd)

	seedPeopleCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.AddCommand(seedPeopleCmd)

	seedProjectsCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.AddCommand(seedProjectsCmd)

	seedCandidateRecordsCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.AddCommand(seedCandidateRecordsCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed",
}

var seedOrganizationsCmd = &cobra.Command{
	Use:   "organizations",
	Short: "seed organization data",
	RunE: func(cmd *cobra.Command, args []string) error {
		authorityClient, err := authority.New(authority.Config{
			MongoDBURI: config.MongoDBURL,
			ESURI:      []string{config.Frontend.Es6URL},
		})
		if err != nil {
			logger.Error("fatal: can't create authority client", "error", err)
			os.Exit(1)
		}

		if err := authorityClient.EnsureOrganizationSeedIndexExists(); err != nil {
			return err
		}

		if force, _ := cmd.Flags().GetBool("force"); !force {
			count, err := authorityClient.CountOrganizations()
			if err != nil {
				return err
			}
			if count > 0 {
				logger.Warn("not seeding dummy data because the database is not empty")
				return nil
			}
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if err := authorityClient.SeedOrganization(scanner.Bytes()); err != nil {
				return err
			}
			logger.Info("imported organization")
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		return nil
	},
}

var seedPeopleCmd = &cobra.Command{
	Use:   "people",
	Short: "seed person data",
	RunE: func(cmd *cobra.Command, args []string) error {
		authorityClient, err := authority.New(authority.Config{
			MongoDBURI: config.MongoDBURL,
			ESURI:      []string{config.Frontend.Es6URL},
		})
		if err != nil {
			logger.Error("fatal: can't create authority client", "error", err)
			os.Exit(1)
		}

		if err := authorityClient.EnsurePersonSeedIndexExists(); err != nil {
			return err
		}

		if force, _ := cmd.Flags().GetBool("force"); !force {
			count, err := authorityClient.CountPeople()
			if err != nil {
				return err
			}
			if count > 0 {
				logger.Warn("not seeding dummy data because the database is not empty")
				return nil
			}
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if err := authorityClient.SeedPerson(scanner.Bytes()); err != nil {
				return err
			}
			logger.Info("imported person")
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		return nil
	},
}

var seedProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "seed project data",
	RunE: func(cmd *cobra.Command, args []string) error {
		authorityClient, err := authority.New(authority.Config{
			MongoDBURI: config.MongoDBURL,
			ESURI:      []string{config.Frontend.Es6URL},
		})
		if err != nil {
			logger.Error("fatal: can't create authority client", "error", err)
			os.Exit(1)
		}

		if err := authorityClient.EnsureProjectSeedIndexExists(); err != nil {
			return err
		}

		if force, _ := cmd.Flags().GetBool("force"); !force {
			count, err := authorityClient.CountProjects()
			if err != nil {
				return err
			}
			if count > 0 {
				logger.Warn("not seeding dummy data because the database is not empty")
				return nil
			}
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if err := authorityClient.SeedProject(scanner.Bytes()); err != nil {
				return err
			}
			logger.Info("imported project")
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		return nil
	},
}

var seedCandidateRecordsCmd = &cobra.Command{
	Use:   "candidate-records",
	Short: "seed candidate record data",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := newServices().Repo

		if force, _ := cmd.Flags().GetBool("force"); !force {
			exists, err := repo.HasCandidateRecords(cmd.Context())
			if err != nil {
				return err
			}
			if exists {
				logger.Warn("not seeding dummy data because the database is not empty")
				return nil
			}
		}

		dec := json.NewDecoder(os.Stdin)
		for {
			rec := &models.CandidateRecord{}
			if err := dec.Decode(rec); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				return err
			}
			if err := repo.AddCandidateRecord(context.TODO(), rec); err != nil {
				return err
			}
			logger.Info("imported candidate record", "sourceID", rec.SourceID, "source", rec.SourceName)
		}

		return nil
	},
}
