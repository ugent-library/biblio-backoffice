package mutate

import (
	"testing"

	"github.com/ugent-library/biblio-backoffice/models"
)

func TestSetPublisher(t *testing.T) {
	p := &models.Publication{}

	publisher := "my-publisher"
	SetPublisher(p, []string{publisher})

	if p.Publisher != publisher {
		t.Errorf("expected publication.Publisher to be '%s', got '%s'", publisher, p.Publisher)
	}
}

func TestRemovePublisher(t *testing.T) {
	p := &models.Publication{Publisher: "my-publisher"}

	RemovePublisher(p, nil)

	if p.Publisher != "" {
		t.Errorf("expected publication.Publisher to be '', got '%s'", p.Publisher)
	}
}
