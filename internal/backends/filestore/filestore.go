package filestore

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path"
)

type Store struct {
	dir     string
	tempDir string
}

func New(basePath string) (*Store, error) {
	s := &Store{
		dir:     path.Join(basePath, "root"),
		tempDir: path.Join(basePath, "tmp"),
	}
	if err := os.MkdirAll(s.dir, os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.tempDir, os.ModePerm); err != nil {
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

func (s *Store) relativeFilePath(checksum string) string {
	fnv32 := fmt.Sprintf("%d", fnvHash(checksum))
	return path.Join(segmentedPath(fnv32, 3), checksum)
}

func (s *Store) filePath(checksum string) string {
	return path.Join(s.dir, s.relativeFilePath(checksum))
}

func (s *Store) Exists(ctx context.Context, checksum string) (bool, error) {
	_, err := os.Stat(s.filePath(checksum))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *Store) Get(ctx context.Context, checksum string) (io.ReadCloser, error) {
	p := s.filePath(checksum)
	if _, err := os.Stat(p); err != nil {
		return nil, err
	}
	return os.Open(p)
}

func (s *Store) Add(ctx context.Context, r io.Reader, oldChecksum string) (string, error) {
	tmpFile, err := os.CreateTemp(s.tempDir, "")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	hasher := sha256.New()

	w := io.MultiWriter(tmpFile, hasher)

	bytesWritten, err := io.Copy(w, r)
	if err != nil {
		return "", err
	}
	if bytesWritten == 0 {
		return "", errors.New("file can't be empty")
	}

	checksum := fmt.Sprintf("%x", hasher.Sum(nil))

	//check sha256 if given
	if oldChecksum != "" && oldChecksum != checksum {
		return "", fmt.Errorf(
			"sha256 checksum did not match '%s', got '%s'",
			oldChecksum,
			checksum,
		)
	}

	//write to final location
	fnv32 := fmt.Sprintf("%d", fnvHash(checksum))
	segmentedPath := segmentedPath(fnv32, 3)
	pathToDir := path.Join(s.dir, segmentedPath)
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

// TODO remove empty intermediate directories?
func (s *Store) Delete(ctx context.Context, checksum string) error {
	return os.Remove(s.filePath(checksum))
}

func (s *Store) DeleteAll(ctx context.Context) error {
	return os.RemoveAll(s.dir)
}
