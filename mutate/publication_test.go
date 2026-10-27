package mutate

import (
	"strings"
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

	doi := "10.52521/kg.v23i1.16846 "
	if err := SetDOI(p, []string{doi}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantDOI := strings.TrimSpace(doi)
	if p.DOI != wantDOI {
		t.Errorf("expected publication.DOI to be '%s', got '%s'", wantDOI, p.DOI)
	}
}

func TestRemoveDOI(t *testing.T) {
	p := &models.Publication{DOI: "10.52521/kg.v23i1.16846 "}

	if err := RemoveDOI(p, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.DOI != "" {
		t.Errorf("expected publication.DOI to be '', got '%s'", p.DOI)
	}
}
