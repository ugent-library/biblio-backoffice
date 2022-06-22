package mutantdb_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/ugent-library/biblio-backend/internal/mutantdb"
	"golang.org/x/net/context"
)

var defaultPgURL = "postgres://localhost:5432/test_mutantdb?sslmode=disable"

type Book struct {
	Title  string
	Author []string
}

type Author struct {
	Name string
}

func newBook() Book      { return Book{} }
func newAuthor() *Author { return &Author{} }

func validateBook(b Book) error {
	if b.Title == "" {
		return errors.New("title missing")
	}
	return nil
}

func TestMutantDB(t *testing.T) {
	ctx := context.Background()

	bookType := mutantdb.NewType("Book", newBook).WithValidator(validateBook)
	authorType := mutantdb.NewType("Author", newAuthor)

	titleAdder := mutantdb.NewMutator("SetTitle", func(b Book, v string) (Book, error) {
		b.Title = v
		return b, nil
	})
	authorAdder := mutantdb.NewMutator("AddAuthor", func(b Book, v string) (Book, error) {
		b.Author = append(b.Author, v)
		return b, nil
	})

	nameAdder := mutantdb.NewMutator("SetName", func(a *Author, v string) (*Author, error) {
		a.Name = v
		return a, nil
	})

	idGenerator := func() (string, error) {
		id := uuid.NewString()
		return "custom:" + id, nil
	}

	pgURL, ok := os.LookupEnv("PG_URL")
	if !ok {
		pgURL = defaultPgURL
	}

	conn, err := pgx.Connect(ctx, pgURL)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(ctx)

	bookStore := mutantdb.NewStore(conn, bookType).
		WithMutators(
			titleAdder,
			authorAdder,
		)
	authorStore := mutantdb.NewStore(conn, authorType).
		WithIDGenerator(idGenerator).
		WithMutators(nameAdder)

	bookID := uuid.NewString()
	bookTitle := "My Title"
	authorID := "author1"
	authorName := "Mr Smith"

	// append
	bookProjectionBeforeTx, err := bookStore.Append(ctx, bookID,
		titleAdder.New(bookTitle),
	)
	if err != nil {
		t.Errorf("got error, want nil: %s", err)
	}

	if err == nil && bookProjectionBeforeTx.Data.Title != bookTitle {
		t.Errorf("got %s, want %s", bookProjectionBeforeTx.Data.Title, bookTitle)
	}

	// rollback transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	if _, err = authorStore.Tx(tx).Append(ctx, authorID, nameAdder.New(authorName)); err != nil {
		t.Errorf("got error, want nil: %s", err)
	}
	if _, err = bookStore.Tx(tx).Append(ctx, bookID, authorAdder.New(authorID)); err != nil {
		t.Errorf("got error, want nil: %s", err)
	}

	if err = tx.Rollback(ctx); err != nil {
		t.Fatal(err)
	}

	bookProjectionAfterTx, err := bookStore.Get(ctx, bookID)
	if err != nil {
		t.Errorf("got error, want nil: %s", err)
	}
	if bookProjectionAfterTx.MutationID != bookProjectionBeforeTx.MutationID {
		t.Errorf("rollback failed")
	}

	// transaction
	tx, err = conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	if _, err = authorStore.Tx(tx).Append(ctx, authorID, nameAdder.New(authorName)); err != nil {
		t.Errorf("got error, want nil: %s", err)
	}
	if _, err = bookStore.Tx(tx).Append(ctx, bookID, authorAdder.New(authorID)); err != nil {
		t.Errorf("got error, want nil: %s", err)
	}

	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	bookProjectionAfterTx, err = bookStore.Get(ctx, bookID)
	if err != nil {
		t.Errorf("got error, want nil: %s", err)
	}
	if bookProjectionAfterTx.MutationID == bookProjectionBeforeTx.MutationID {
		t.Errorf("commit failed")
	}

	// test custom id generator
	authorProjection, err := authorStore.Get(ctx, authorID)
	if err != nil {
		t.Errorf("got error, want nil: %s", err)
	}
	if !strings.HasPrefix(authorProjection.MutationID, "custom:") {
		t.Error("got default id, want custom id")
	}

	// test GetAll
	var authors []*Author

	c, err := authorStore.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	for c.HasNext() {
		p, err := c.Next()
		if err != nil {
			t.Errorf("got error, want nil: %s", err)
		}
		authors = append(authors, p.Data)
	}
	if err := c.Error(); err != nil {
		t.Errorf("got error, want nil: %s", err)
	}
	if len(authors) != 1 {
		t.Errorf("got %d, want 1", len(authors))
	} else if authors[0].Name != authorName {
		t.Errorf("got %s, want %s", authors[0].Name, authorName)
	}

	// test GetAt
	oldBookProjection, err := bookStore.GetAt(ctx, bookID, bookProjectionBeforeTx.MutationID)
	if err != nil {
		t.Fatal(err)
	}
	if len(oldBookProjection.Data.Author) != len(bookProjectionBeforeTx.Data.Author) {
		t.Errorf("got %d, want %d", len(oldBookProjection.Data.Author), len(bookProjectionBeforeTx.Data.Author))
	}

	// test conflict detection
	_, err = bookStore.AppendAfter(ctx, bookID, oldBookProjection.MutationID,
		titleAdder.New(bookTitle),
	)
	var conflict *mutantdb.ErrConflict
	if err == nil {
		t.Error("conflict detection failed")
	}
	if !errors.As(err, &conflict) {
		t.Error("conflict detection failed")
	}
}
