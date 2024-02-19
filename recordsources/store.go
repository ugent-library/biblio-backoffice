package recordsources

import (
	"context"
	"io"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/backends"
)

func StoreURL(ctx context.Context, url string, f backends.FileStore) (string, int, error) {
	// TODO timeouts
	res, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()
	cr := &countingReader{r: res.Body}

	sha256, err := f.Add(ctx, cr, "")

	return sha256, cr.n, err
}

type countingReader struct {
	r io.Reader
	n int
}

func (r *countingReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.n += n
	return n, err
}
