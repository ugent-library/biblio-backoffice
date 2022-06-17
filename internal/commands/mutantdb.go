package commands

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/mutantdb"
	"github.com/ugent-library/biblio-backend/internal/ulid"
)

func init() {
	rootCmd.AddCommand(testEventstoreCmd)
}

var DatasetType = mutantdb.NewType("Dataset", NewDataset).WithValidator(func(d *models.Dataset) error { return d.Validate() })
var PublicationType = mutantdb.NewType("Publication", NewPublication)

var DatasetReplacer = mutantdb.NewMutator("Replace", ReplaceDataset)
var DatasetAbstractAdder = mutantdb.NewMutator("AddAbstract", AddDatasetAbstract)
var DatasetPublicationAdder = mutantdb.NewMutator("AddPublication", AddDatasetPublication)

var PublicationDatasetAdder = mutantdb.NewMutator("AddDataset", AddPublicationDataset)

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
		ctx := context.Background()

		conn, err := pgxpool.Connect(ctx, viper.GetString("pg-conn"))
		if err != nil {
			log.Fatal(err)
		}

		datasetStore := mutantdb.NewStore(conn, DatasetType).
			WithIDGenerator(ulid.Generate).
			WithMutators(
				DatasetReplacer,
				DatasetAbstractAdder,
				DatasetPublicationAdder,
			)
		pubStore := mutantdb.NewStore(conn, PublicationType).
			WithIDGenerator(ulid.Generate).
			WithMutators(
				PublicationDatasetAdder,
			)

		// TEST Append
		datasetID := ulid.MustGenerate()
		pubID := ulid.MustGenerate()

		_, err = datasetStore.Append(ctx, datasetID,
			DatasetReplacer.New(
				&models.Dataset{Status: "private", Title: "Test dataset"},
			),
			DatasetAbstractAdder.New(
				models.Text{Lang: "eng", Text: "Test abstract"},
				mutantdb.Meta{"UserID": "123"},
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		pBeforeTx, err := datasetStore.Get(ctx, datasetID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", pBeforeTx.Data)

		// TEST TRANSACTIONS

		// rollback
		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback(ctx)

		_, err = datasetStore.Tx(tx).Append(ctx, datasetID, DatasetPublicationAdder.New(pubID))
		if err != nil {
			log.Fatal(err)
		}
		_, err = pubStore.Tx(tx).Append(ctx, pubID, PublicationDatasetAdder.New(datasetID))
		if err != nil {
			log.Fatal(err)
		}

		if err = tx.Rollback(ctx); err != nil {
			log.Fatal(err)
		}

		pAfterTx, err := datasetStore.Get(ctx, datasetID)
		if pAfterTx.MutationID != pBeforeTx.MutationID {
			log.Fatalf("Rollback failed, mutation id changed from %s, to %s", pBeforeTx.MutationID, pAfterTx.MutationID)
		}

		// success
		tx, err = conn.Begin(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback(ctx)

		_, err = datasetStore.Tx(tx).Append(ctx, datasetID, DatasetPublicationAdder.New(pubID))
		if err != nil {
			log.Fatal(err)
		}
		_, err = pubStore.Tx(tx).Append(ctx, pubID, PublicationDatasetAdder.New(datasetID))
		if err != nil {
			log.Fatal(err)
		}

		if err = tx.Commit(ctx); err != nil {
			log.Fatal(err)
		}

		// TEST Get

		p, err := datasetStore.Get(ctx, datasetID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("projection data after tx: %+v", p.Data)

		// TEST GetAll

		c, err := datasetStore.GetAll(ctx)
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

		// TEST GetAt

		p, err = datasetStore.GetAt(ctx, datasetID, pBeforeTx.MutationID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("projection before tx: %+v", p)
		log.Printf("projection data before tx: %+v", p.Data)

		// TEST conflict detection
		_, err = datasetStore.AppendAfter(ctx, datasetID, pBeforeTx.MutationID,
			DatasetAbstractAdder.New(models.Text{Lang: "eng", Text: "Test abstract"}),
		)
		if err == nil {
			log.Fatal("conflict detection failed")
		} else {
			log.Printf("invalid AfterMutation gives conflict error: %s", err)
		}

		// TEST Append & Get

		anyP, err := datasetStore.Append(ctx, datasetID,
			DatasetAbstractAdder.New(models.Text{Lang: "eng", Text: "Another test abstract"}),
		)
		if err != nil {
			log.Printf("invalid AfterMutation gives conflict error: %s", err)
		}
		log.Printf("get projection after append: %+v", anyP)
		log.Printf("get projection data after append: %+v", anyP.Data)
	},
}
