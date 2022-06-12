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

var SetHandler = eventstore.NewEventHandler("Set", SetDataset)

var AddAbstractHandler = eventstore.NewEventHandler("AddAbstract", AddAbstract)

func SetDataset(data *models.Dataset, newData *models.Dataset) (*models.Dataset, error) {
	return newData, nil
}

func AddAbstract(data *models.Dataset, a models.Text) (*models.Dataset, error) {
	data.Abstract = append(data.Abstract, a)
	return data, nil
}

var testEventstoreCmd = &cobra.Command{
	Use: "test-eventstore",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := eventstore.Connect(context.Background(), viper.GetString("pg-conn"),
			eventstore.WithIDGenerator(ulid.Generate),
		)
		if err != nil {
			log.Fatal(err)
		}

		datasetRepository := eventstore.NewRepository[*models.Dataset](store, "Dataset")
		// datasetRepository.AddEventHandlers(
		// 	SetHandler,
		// 	AddAbstractHandler,
		// )

		streamID := ulid.MustGenerate()

		err = store.Append(context.Background(),
			SetHandler.NewEvent(
				datasetRepository.StreamType(),
				streamID,
				&models.Dataset{Title: "Test dataset", Publisher: "Test publisher"},
			),
			AddAbstractHandler.NewEvent(
				datasetRepository.StreamType(),
				streamID,
				models.Text{Lang: "eng", Text: "Test abstract"},
				map[string]string{"UserID": "123"},
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		p, err := datasetRepository.Get(context.Background(), streamID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", p)
		log.Printf("%+v", p.Data)

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
	},
}
