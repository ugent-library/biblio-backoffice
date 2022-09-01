package server

import (
	"bufio"
	"fmt"
	"io"
	"os"

	api "github.com/ugent-library/biblio-backend/api/v1"
)

const fileBufSize = 524288

func (s *server) GetFile(req *api.GetFileRequest, stream api.Biblio_GetFileServer) error {
	fPath := s.services.FileStore.FilePath(req.Sha256)
	f, err := os.Open(fPath)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	buf := make([]byte, fileBufSize)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("cannot read chunk to buffer: %w", err)
		}

		req := &api.GetFileResponse{Chunk: buf[:n]}

		if err := stream.Send(req); err != nil {
			return fmt.Errorf("cannot send chunk to client: %w", err)
		}
	}

	return nil
}

func (s *server) AddFile(stream api.Biblio_AddFileServer) error {
	var (
		sha256       string
		fileStoreErr error
	)

	pr, pw := io.Pipe()
	waitc := make(chan struct{})

	go func() {
		sha256, fileStoreErr = s.services.FileStore.Add(pr)
		close(waitc)
	}()

	// TODO break if fileStore add returns an error (use select)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if _, err := pw.Write(req.Chunk); err != nil {
			return err
		}
	}

	pw.Close()

	<-waitc

	if fileStoreErr != nil {
		return fileStoreErr
	}
	return stream.SendAndClose(&api.AddFileResponse{Sha256: sha256})
}
