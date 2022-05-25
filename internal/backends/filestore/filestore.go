package filestore

import (
	"crypto/sha256"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type Store struct {
	rootPath string
	tmpPath  string
}

func New(basePath string) (*Store, error) {
	s := &Store{
		rootPath: path.Join(basePath, "root"),
		tmpPath:  path.Join(basePath, "tmp"),
	}
	if err := os.MkdirAll(s.rootPath, os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.tmpPath, os.ModePerm); err != nil {
		return nil, err
	}
	return s, nil
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

func (s *Store) RelativeFilePath(checksum string) string {
	fnv32 := fmt.Sprintf("%d", fnvHash(checksum))
	return path.Join(segmentedPath(fnv32, 3), checksum)
}

func (s *Store) FilePath(checksum string) string {
	return path.Join(s.rootPath, s.RelativeFilePath(checksum))
}

func (s *Store) Add(r io.Reader) (string, error) {
	tmpFile, err := ioutil.TempFile(s.tmpPath, "")
	if err != nil {
		return "", err
	}
	// fmt.Printf("%s\n", tmpFile.Name())
	defer os.Remove(tmpFile.Name())

	hash := sha256.New()

	w := io.MultiWriter(tmpFile, hash)

	if _, err := io.Copy(w, r); err != nil {
		return "", err
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	fnv32 := fmt.Sprintf("%d", fnvHash(checksum))

	// log.Printf("sha256: %s", checksum)
	// log.Printf("fnv: %s", fnv32)

	segmentedPath := segmentedPath(fnv32, 3)
	// log.Printf("segmented path: %s", segmentedPath)
	pathToDir := path.Join(s.rootPath, segmentedPath)
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

// TODO remove empty intermediate directories
func (s *Store) Purge(checksum string) error {
	return os.Remove(s.FilePath(checksum))
}

func (s *Store) PurgeAll() error {
	dir, err := os.Open(s.rootPath)
	if err != nil {
		return err
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		if err := os.RemoveAll(path.Join(s.rootPath, name)); err != nil {
			return err
		}
	}
	return nil
}
