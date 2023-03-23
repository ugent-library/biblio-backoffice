package server

import (
	"bufio"
	"fmt"
	"io"
	"os"

	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

		req := &api.GetFileResponse{
			Chunk: buf[:n],
		}

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

recv:
	for {
		select {
		case <-waitc:
			break recv
		default:
			req, err := stream.Recv()
			if err == io.EOF {
				break recv
			}
			if err != nil {
				return status.Errorf(codes.Internal, "failed to read stream: %s", err)
			}
			if _, err := pw.Write(req.Chunk); err != nil {
				return status.Errorf(codes.Internal, "failed to write file chunk: %s", err)
			}
		}
	}

	pw.Close()
	<-waitc

	if fileStoreErr != nil {
		return status.Errorf(codes.Internal, "failed to write to stream: %v", fileStoreErr)
	}

	if err := stream.SendAndClose(&api.AddFileResponse{
		Response: &api.AddFileResponse_Sha256{
			Sha256: sha256,
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to write to stream: %v", err)
	}

	return nil
}
