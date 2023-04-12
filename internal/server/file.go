package server

import (
	"bufio"
	"context"
	"fmt"
	"io"

	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const fileBufSize = 524288

func (s *server) GetFile(req *api.GetFileRequest, stream api.Biblio_GetFileServer) error {
	b, err := s.services.FileStore.Get(stream.Context(), req.Sha256)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer b.Close()

	r := bufio.NewReader(b)
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

func (s *server) ExistsFile(ctx context.Context, req *api.ExistsFileRequest) (*api.ExistsFileResponse, error) {
	exists, err := s.services.FileStore.Exists(ctx, req.Sha256)
	return &api.ExistsFileResponse{
		Exists: exists,
	}, err
}

func (s *server) AddFile(stream api.Biblio_AddFileServer) error {
	var (
		sha256       string
		fileStoreErr error
	)

	pr, pw := io.Pipe()
	waitc := make(chan struct{})

	go func() {
		sha256, fileStoreErr = s.services.FileStore.Add(stream.Context(), pr, "")
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
