package commands

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/mutantdb"
	"github.com/ugent-library/biblio-backend/internal/ulid"
)

func init() {
	rootCmd.AddCommand(testEventstoreCmd)
}

var DatasetType = mutantdb.NewType("Dataset", NewDataset)
var PublicationType = mutantdb.NewType("Publication", NewPublication)

var DatasetReplacer = mutantdb.NewMutator(DatasetType, "Replace", ReplaceDataset)
var DatasetAbstractAdder = mutantdb.NewMutator(DatasetType, "AddAbstract", AddDatasetAbstract)
var DatasetPublicationAdder = mutantdb.NewMutator(DatasetType, "AddPublication", AddDatasetPublication)

var PublicationDatasetAdder = mutantdb.NewMutator(PublicationType, "AddDataset", AddPublicationDataset)

func NewDataset() *models.Dataset {
	return &models.Dataset{
		Status: "private",
	}
}

func ReplaceDataset(data *models.Dataset, newData *models.Dataset) (*models.Dataset, error) {
	return newData, nil
}

func AddDatasetAbstract(data *models.Dataset, a models.Text) (*models.Dataset, error) {
	data.Abstract = append(data.Abstract, a)
	return data, nil
}

func AddDatasetPublication(data *models.Dataset, pubID string) (*models.Dataset, error) {
	data.RelatedPublication = append(data.RelatedPublication, models.RelatedPublication{ID: pubID})
	return data, nil
}

func NewPublication() *models.Publication {
	return &models.Publication{
		Status: "private",
	}
}

func AddPublicationDataset(data *models.Publication, datasetID string) (*models.Publication, error) {
	data.RelatedDataset = append(data.RelatedDataset, models.RelatedDataset{ID: datasetID})
	return data, nil
}

var testEventstoreCmd = &cobra.Command{
	Use: "test-mutantdb",
	Run: func(cmd *cobra.Command, args []string) {
		// TEST STORE
		ctx := context.Background()

		store, err := mutantdb.Connect(ctx, viper.GetString("pg-conn"),
			mutantdb.WithIDGenerator(ulid.Generate),
			mutantdb.WithMutators(
				DatasetReplacer,
				DatasetAbstractAdder,
				DatasetPublicationAdder,
				PublicationDatasetAdder,
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		// TEST Append
		datasetID := ulid.MustGenerate()
		pubID := ulid.MustGenerate()

		err = store.Append(datasetID,
			DatasetReplacer.New(
				&models.Dataset{Title: "Test dataset", Publisher: "Test publisher"},
			),
			DatasetAbstractAdder.New(
				models.Text{Lang: "eng", Text: "Test abstract"},
				mutantdb.Meta{"UserID": "123"},
			),
		).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}

		datasetRepository := mutantdb.NewRepository(store, DatasetType)
		pBeforeTx, err := datasetRepository.Get(ctx, datasetID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", pBeforeTx.Data)

		// TEST TRANSACTIONS
		tx, err := store.BeginTx(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback(ctx)

		if err = tx.Append(datasetID, DatasetPublicationAdder.New(pubID)).Do(ctx); err != nil {
			log.Fatal(err)
		}
		if err = tx.Append(pubID, PublicationDatasetAdder.New(datasetID)).Do(ctx); err != nil {
			log.Fatal(err)
		}

		if err = tx.Commit(ctx); err != nil {
			log.Fatal(err)
		}

		// test repository Get
		p, err := datasetRepository.Get(ctx, datasetID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("projection data after tx: %+v", p.Data)

		// test repository GetAll
		c, err := datasetRepository.GetAll(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()
		for c.HasNext() {
			p, err := c.Next()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("iterated id: %s", p.ID)
		}
		if err := c.Error(); err != nil {
			log.Fatal(err)
		}

		// test repository GetAt
		p, err = datasetRepository.GetAt(ctx, datasetID, pBeforeTx.MutationID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("projection before tx: %+v", p)
		log.Printf("projection data before tx: %+v", p.Data)

		// TEST invalid Append
		err = store.Append(datasetID,
			DatasetAbstractAdder.New(
				models.Text{Lang: "eng", Text: "Test abstract"},
				mutantdb.Meta{"UserID": "123"},
			),
			PublicationDatasetAdder.New(datasetID),
		).Do(ctx)
		if err != nil {
			log.Printf("invalid append gives error: %s", err)
		}
	},
}
