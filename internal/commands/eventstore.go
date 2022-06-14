package commands

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/eventstore"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/ulid"
)

func init() {
	rootCmd.AddCommand(testEventstoreCmd)
}

// TODO use stream type everywhere
var DatasetType = eventstore.NewStreamType("Dataset", NewDataset)

// TODO clearer distinction between event and command
var ReplaceDatasetHandler = eventstore.NewEventHandler(DatasetType, "Replaced", ReplaceDataset)
var AddDatasetAbstractHandler = eventstore.NewEventHandler(DatasetType, "AbstractAdded", AddDatasetAbstract)

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
	Use: "test-eventstore",
	Run: func(cmd *cobra.Command, args []string) {
		// TEST STORE
		store, err := eventstore.Connect(context.Background(), viper.GetString("pg-conn"),
			eventstore.WithIDGenerator(ulid.Generate),
		)
		if err != nil {
			log.Fatal(err)
		}

		store.AddEventHandlers(
			ReplaceDatasetHandler,
			AddDatasetAbstractHandler,
		)

		// test Append
		streamID := ulid.MustGenerate()

		err = store.Append(context.Background(),
			ReplaceDatasetHandler.NewEvent(
				streamID,
				&models.Dataset{Title: "Test dataset", Publisher: "Test publisher"},
			),
			AddDatasetAbstractHandler.NewEvent(
				streamID,
				models.Text{Lang: "eng", Text: "Test abstract"},
				eventstore.Meta{"UserID": "123"},
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		// TEST REPOSITORY
		datasetRepository := eventstore.NewRepository(store, DatasetType)

		// test Get
		p, err := datasetRepository.Get(context.Background(), streamID)
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
			log.Printf("iterated id %s", p.StreamID)
		}
		if err := c.Error(); err != nil {
			log.Fatal(err)
		}

		// test GetAt
		// p, err = datasetRepository.GetAt(context.Background(), "01G5E2D1HYK531S6G48PM9WBW8", "01G5E2D1HYM158TRFZJBAWJ22Q")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// log.Printf("%+v", p)
		// log.Printf("%+v", p.Data)
	},
}
