package engine

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/ugent-library/biblio-backend/internal/models"
)

func (e *Engine) UpdatePublication(p *models.Publication) (*models.Publication, error) {
	p.Vacuum()

	if err := p.Validate(); err != nil {
		log.Printf("%#v", err)
		return nil, err
	}

	p, err := e.StorageService.UpdatePublication(p)
	if err != nil {
		return nil, err
	}

	if err := e.PublicationSearchService.IndexPublication(p); err != nil {
		log.Printf("error indexing publication %+v", err)
		return nil, err
	}

	return p, nil
}

func (e *Engine) Publications(args *models.SearchArgs) (*models.PublicationHits, error) {
	args = args.Clone().WithFilter("status", "private", "public")
	return e.PublicationSearchService.SearchPublications(args)
}

func (e *Engine) UserPublications(userID string, args *models.SearchArgs) (*models.PublicationHits, error) {
	args = args.Clone().WithFilter("status", "private", "public")
	switch args.FilterFor("scope") {
	case "created":
		args.WithFilter("creator_id", userID)
	case "contributed":
		args.WithFilter("author.id", userID)
	default:
		args.WithFilter("creator_id|author.id", userID)
	}
	delete(args.Filters, "scope")
	return e.PublicationSearchService.SearchPublications(args)
}

func (e *Engine) BatchPublishPublications(userID string, args *models.SearchArgs) (err error) {
	var hits *models.PublicationHits
	for {
		hits, err = e.UserPublications(userID, args)
		for _, pub := range hits.Hits {
			pub.Status = "public"
			if _, err = e.UpdatePublication(pub); err != nil {
				break
			}
		}
		if !hits.NextPage() {
			break
		}
		args.Page = args.Page + 1
	}
	return
}

func (e *Engine) GetPublicationDatasets(id string) ([]*models.Dataset, error) {
	return nil, nil
}

func (e *Engine) AddPublicationDataset(id, datasetID string) error {
	return errors.New("not implemented")
}

func (e *Engine) RemovePublicationDataset(id, datasetID string) error {
	return errors.New("not implemented")
}

func (e *Engine) ImportUserPublicationByIdentifier(userID, source, identifier string) (*models.Publication, error) {
	return nil, errors.New("not implemented")
}

func (e *Engine) ImportUserPublications(userID, source string, file io.Reader) (string, error) {
	return "", errors.New("not implemented")
}

func (c *Engine) ServePublicationThumbnail(fileURL string, w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func fnvHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func segmentedPath(str string, size int) string {
	strLength := len(str)
	var segments []string
	var stop int
	for i := 0; i < strLength; i += size {
		stop = i + size
		if stop > strLength {
			stop = strLength
		}
		segments = append(segments, str[i:stop])
	}
	return path.Join(segments...)
}

func (e *Engine) FilePath(checksum string) string {
	fnv32 := fmt.Sprintf("%d", fnvHash(checksum))
	return path.Join("/Users/nsteenla/tmp/biblio_backend/files", segmentedPath(fnv32, 3), checksum)
}

func (e *Engine) StoreFile(r io.Reader) (string, error) {
	tmpFile, err := ioutil.TempFile("/Users/nsteenla/tmp/biblio_backend/tmp", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", tmpFile.Name())
	defer os.Remove(tmpFile.Name())

	hash := sha256.New()

	w := io.MultiWriter(tmpFile, hash)

	if _, err := io.Copy(w, r); err != nil {
		log.Println("error copying file", err)
		return "", err
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	fnv32 := fmt.Sprintf("%d", fnvHash(checksum))

	log.Printf("sha256: %s", checksum)
	log.Printf("fnv: %s", fnv32)
	log.Printf("segmented path: %s", segmentedPath(fnv32, 3))

	pathToDir := path.Join("/Users/nsteenla/tmp/biblio_backend/files", segmentedPath(fnv32, 3))
	pathToFile := path.Join(pathToDir, checksum)

	// file already stored
	if _, err := os.Stat(pathToFile); !os.IsNotExist(err) {
		return checksum, nil
	}

	if err := os.MkdirAll(pathToDir, os.ModePerm); err != nil {
		return "", err
	}

	if err := os.Rename(tmpFile.Name(), path.Join(pathToDir, checksum)); err != nil {
		return "", err
	}

	return checksum, nil
}
