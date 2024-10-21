// TODO use fs.FS
package fsstore

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

type Config struct {
	Dir     string
	TempDir string
}

type Store struct {
	dir, tempDir string
}

func New(c Config) (*Store, error) {
	s := &Store{
		dir:     c.Dir,
		tempDir: c.TempDir,
	}
	if err := checkDirOk(s.dir); err != nil {
		return nil, fmt.Errorf("fsstore: store dir: %w", err)
	}
	if err := checkDirOk(s.tempDir); err != nil {
		return nil, fmt.Errorf("fsstore: temp dir: %w", err)
	}
	return s, nil
}

func checkDirOk(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("can't stat path %s: %w", dir, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("path %s is not a directory", dir)
	}
	return nil
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

func (s *Store) filePath(checksum string) string {
	fnv32 := fmt.Sprintf("%d", fnvHash(checksum))
	return path.Join(s.dir, segmentedPath(fnv32, 3), checksum)
}

func (s *Store) Exists(ctx context.Context, checksum string) (bool, error) {
	fp := s.filePath(checksum)
	_, err := os.Stat(fp)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("fsstore.Exists: can't stat file %s for checksum %s: %w", fp, checksum, err)
}

func (s *Store) Get(ctx context.Context, checksum string) (io.ReadCloser, error) {
	fp := s.filePath(checksum)
	r, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("fsstore.Get: can't open file %s for checksum %s: %w", fp, checksum, err)
	}
	return r, nil
}

func (s *Store) Add(ctx context.Context, r io.Reader, oldChecksum string) (string, error) {
	tmpFile, err := os.CreateTemp(s.tempDir, "")
	if err != nil {
		return "", fmt.Errorf("fsstore.Add: can't create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	hasher := sha256.New()

	w := io.MultiWriter(tmpFile, hasher)

	bytesWritten, err := io.Copy(w, r)
	if err != nil {
		return "", fmt.Errorf("fsstore.Add: write failed: %w", err)
	}
	if bytesWritten == 0 {
		return "", errors.New("fsstore.Add: file can't be empty")
	}

	checksum := fmt.Sprintf("%x", hasher.Sum(nil))

	// check sha256 if given
	if oldChecksum != "" && oldChecksum != checksum {
		return "", fmt.Errorf("fsstore.Add: sha256 checksums don't match, expected %q, got %q", oldChecksum, checksum)
	}

	// file already stored
	exists, err := s.Exists(ctx, checksum)
	if err != nil {
		return "", fmt.Errorf("fsstore.Add: %w", err)
	}
	if exists {
		return checksum, nil
	}

	// write to final location
	fnv32 := fmt.Sprintf("%d", fnvHash(checksum))
	fileDirPath := path.Join(s.dir, segmentedPath(fnv32, 3))

	if err := os.MkdirAll(fileDirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("fsstore.Add: can't create dir %s for file with checksum %s: %w", fileDirPath, checksum, err)
	}

	filePath := path.Join(fileDirPath, checksum)

	if err := os.Rename(tmpFile.Name(), filePath); err != nil {
		return "", fmt.Errorf("fsstore.Add: can't move file with checksum %s to %s: %w", checksum, filePath, err)
	}

	return checksum, nil
}

// TODO remove empty intermediate directories?
func (s *Store) Delete(ctx context.Context, checksum string) error {
	fp := s.filePath(checksum)
	if err := os.Remove(fp); err != nil {
		return fmt.Errorf("fsstore.Delete: can't remove file %s for checksum %s: %w", fp, checksum, err)
	}
	return nil
}

func (s *Store) DeleteAll(ctx context.Context) error {
	if err := os.RemoveAll(s.dir); err != nil {
		return fmt.Errorf("fsstore.DeleteAll: can't remove all files: %w", err)
	}
	return nil
}
