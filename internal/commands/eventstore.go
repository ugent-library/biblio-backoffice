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

var testEventstoreCmd = &cobra.Command{
	Use: "test-eventstore",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := eventstore.Connect(context.Background(), viper.GetString("pg-conn"))
		if err != nil {
			log.Fatal(err)
		}

		datasetProcessor := eventstore.NewProcessor[*models.Dataset]()

		datasetProcessor.AddHandler("Set", eventstore.Handler(func(data *models.Dataset, eventData *models.Dataset) (*models.Dataset, error) {
			return eventData, nil
		}))
		datasetProcessor.AddHandler("AddAbstract", eventstore.Handler(func(data *models.Dataset, eventData *models.Text) (*models.Dataset, error) {
			data.Abstract = append(data.Abstract, *eventData)
			return data, nil
		}))

		store.AddProcessor("Dataset", datasetProcessor)

		datasetRepository := eventstore.NewRepository[*models.Dataset](store, "Dataset")

		streamID := ulid.MustGenerate()

		err = store.Append(context.Background(),
			eventstore.Event{
				ID:         ulid.MustGenerate(),
				StreamID:   streamID,
				StreamType: "Dataset",
				Type:       "Set",
				Data:       &models.Dataset{Title: "Test dataset", Publisher: "Test publisher"},
				Meta: map[string]string{
					"UserID": "123",
				},
			},
			eventstore.Event{
				ID:         ulid.MustGenerate(),
				StreamID:   streamID,
				StreamType: "Dataset",
				Type:       "AddAbstract",
				Data:       &models.Text{Lang: "eng", Text: "Test abstract"},
				Meta: map[string]string{
					"UserID": "123",
				},
			},
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
	},
}
