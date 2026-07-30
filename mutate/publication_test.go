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

func TestSetDOI(t *testing.T) {
	p := &models.Publication{}

	doiInput := "10.52521/kg.v23i1.16846 "
	expectedDOI := "10.52521/kg.v23i1.16846"
	if err := SetDOI(p, []string{doiInput}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.DOI != expectedDOI {
		t.Errorf("expected publication.DOI to be '%s', got '%s'", expectedDOI, p.DOI)
	}
}

func TestRemoveDOI(t *testing.T) {
	p := &models.Publication{DOI: "10.52521/kg.v23i1.16846 "}

	RemoveDOI(p, nil)

	if p.DOI != "" {
		t.Errorf("expected publication.DOI to be '', got '%s'", p.DOI)
	}
}
