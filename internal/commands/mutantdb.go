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

var DatasetReplacer = mutantdb.NewMutator(DatasetType, "Replace", ReplaceDataset)
var DatasetAbstractAdder = mutantdb.NewMutator(DatasetType, "AddAbstract", AddDatasetAbstract)

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

var testEventstoreCmd = &cobra.Command{
	Use: "test-mutantdb",
	Run: func(cmd *cobra.Command, args []string) {
		// TEST STORE
		store, err := mutantdb.Connect(context.Background(), viper.GetString("pg-conn"),
			mutantdb.WithIDGenerator(ulid.Generate),
			mutantdb.WithMutators(
				DatasetReplacer,
				DatasetAbstractAdder,
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		// test Append
		entityID := ulid.MustGenerate()

		err = store.Append(context.Background(),
			DatasetReplacer.New(
				entityID,
				&models.Dataset{Title: "Test dataset", Publisher: "Test publisher"},
			),
			DatasetAbstractAdder.New(
				entityID,
				models.Text{Lang: "eng", Text: "Test abstract"},
				mutantdb.Meta{"UserID": "123"},
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		// TEST REPOSITORY
		datasetRepository := mutantdb.NewRepository(store, DatasetType)

		// test Get
		p, err := datasetRepository.Get(context.Background(), entityID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", p)
		log.Printf("%+v", p.Data)

		// test GetAll
		c, err := datasetRepository.GetAll(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()
		for c.HasNext() {
			p, err := c.Next()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("iterated id %s", p.ID)
		}
		if err := c.Error(); err != nil {
			log.Fatal(err)
		}

		// test GetAt
		p, err = datasetRepository.GetAt(context.Background(), "01G5E2D1HYK531S6G48PM9WBW8", "01G5E2D1HYM158TRFZJBAWJ22Q")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", p)
		log.Printf("%+v", p.Data)
	},
}
