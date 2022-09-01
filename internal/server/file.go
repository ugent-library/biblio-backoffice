package server

import (
	"bufio"
	"fmt"
	"io"
	"os"

	api "github.com/ugent-library/biblio-backend/api/v1"
)

func (s *server) GetFile(req *api.GetFileRequest, stream api.Biblio_GetFileServer) error {
	fPath := s.services.FileStore.FilePath(req.Sha256)
	f, err := os.Open(fPath)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	buf := make([]byte, 1024)

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
