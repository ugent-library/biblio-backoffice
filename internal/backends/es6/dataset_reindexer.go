package es6

import (
	"fmt"
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type DatasetReindexer struct {
	reindexer
	Datasets
}

func NewDatasetReindexer(datasets *Datasets) *DatasetReindexer {
	newIndexName := fmt.Sprintf(
		"%s_version_%s",
		datasets.Client.Index,
		time.Now().UTC().Format("20060102150405"),
	)
	alias := datasets.Client.Index
	d := datasets.Clone()
	d.Client.Index = newIndexName

	return &DatasetReindexer{
		reindexer: reindexer{
			Client: d.Client,
			alias:  alias,
		},
		Datasets: *d,
	}
}

func (dr *DatasetReindexer) Reindex(ch <-chan *models.Dataset) error {

	// there must be an alias and an old index
	aliasStatus, aliasStatusErr := dr.reindexer.getAliasStatus(dr.alias)
	if aliasStatusErr != nil {
		return aliasStatusErr
	}
	if len(aliasStatus.Indexes) > 1 {
		return fmt.Errorf("alias %s points to more than one index", aliasStatus.Name)
	}

	oldIndexName := aliasStatus.Indexes[0]

	// create new index (without alias)
	if err := dr.reindexer.createIndex(dr.reindexer.Client.Index, dr.reindexer.Client.Settings); err != nil {
		return err
	}

	// reindex from channel (TODO: wait group?)
	dr.Datasets.IndexMultiple(ch)

	// switch alias from old index to new index
	if err := dr.reindexer.switchAlias(dr.alias, oldIndexName, dr.reindexer.Client.Index); err != nil {
		return err
	}

	// remove old index
	if err := dr.reindexer.deleteIndex(oldIndexName); err != nil {
		return err
	}

	return nil
}
