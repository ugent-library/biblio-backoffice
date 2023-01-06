package es6

import (
	"fmt"
	"time"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type PublicationReindexer struct {
	reindexer
	Publications
}

func NewPublicationReindexer(publications *Publications) *PublicationReindexer {
	newIndexName := fmt.Sprintf(
		"%s_version_%s",
		publications.Client.Index,
		time.Now().UTC().Format("20060102150405"),
	)
	alias := publications.Client.Index
	p := publications.Clone()
	p.Client.Index = newIndexName

	return &PublicationReindexer{
		reindexer: reindexer{
			Client: p.Client,
			alias:  alias,
		},
		Publications: *p,
	}
}

func (pr *PublicationReindexer) Reindex(ch <-chan *models.Publication) error {

	// there must be an alias and an old index
	aliasStatus, aliasStatusErr := pr.reindexer.getAliasStatus(pr.alias)
	if aliasStatusErr != nil {
		return aliasStatusErr
	}
	if len(aliasStatus.Indexes) > 1 {
		return fmt.Errorf("alias %s points to more than one index", aliasStatus.Name)
	}

	oldIndexName := aliasStatus.Indexes[0]

	// create new index (without alias)
	if err := pr.reindexer.createIndex(pr.reindexer.Client.Index, pr.reindexer.Client.Settings); err != nil {
		return err
	}

	// reindex from channel
	pr.Publications.IndexMultiple(ch)

	// switch alias from old index to new index
	if err := pr.reindexer.switchAlias(pr.alias, oldIndexName, pr.reindexer.Client.Index); err != nil {
		return err
	}

	// remove old index
	if err := pr.reindexer.deleteIndex(oldIndexName); err != nil {
		return err
	}

	return nil
}
