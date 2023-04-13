package s3store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type Config struct {
	Endpoint   string
	Region     string
	ID         string
	Secret     string
	Bucket     string
	TempBucket string
}

type Store struct {
	client     *s3.Client
	bucket     string
	tempBucket string
}

func New(c Config) (*Store, error) {
	if c.TempBucket == "" {
		c.TempBucket = c.Bucket
	}

	config := aws.Config{
		Region:      c.Region,
		Credentials: credentials.NewStaticCredentialsProvider(c.ID, c.Secret, ""),
	}
	if c.Endpoint != "" {
		config.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, opts ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               c.Endpoint,
				SigningRegion:     c.Region,
				HostnameImmutable: true,
			}, nil
		})
	}

	return &Store{
		client:     s3.NewFromConfig(config),
		bucket:     c.Bucket,
		tempBucket: c.TempBucket,
	}, nil
}

func (s *Store) Get(ctx context.Context, checksum string) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (s *Store) Exists(ctx context.Context, checksum string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(checksum),
	})
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Store) Add(ctx context.Context, b io.Reader, oldChecksum string) (string, error) {
	tempKey := uuid.New().String()

	hasher := sha256.New()

	tee := io.TeeReader(b, hasher)

	tempExpires := time.Now().Add(time.Hour)

	uploader := manager.NewUploader(s.client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:  aws.String(s.tempBucket),
		Key:     aws.String(tempKey),
		Body:    tee,
		Expires: &tempExpires,
	})
	if err != nil {
		return "", err
	}

	checksum := fmt.Sprintf("%x", hasher.Sum(nil))

	// check sha256 if given
	if oldChecksum != "" && oldChecksum != checksum {
		return "", fmt.Errorf("sha256 checksum did not match '%s', got '%s'", oldChecksum, checksum)
	}

	_, err = s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.bucket),
		CopySource: aws.String(s.tempBucket + "/" + tempKey),
		Key:        aws.String(s.bucket + "/" + checksum),
	})
	if err != nil {
		return "", err
	}

	return checksum, nil
}

func (s *Store) Delete(ctx context.Context, checksum string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(checksum),
	})
	return err
}

func (s *Store) DeleteAll(ctx context.Context) error {
	return errors.New("not implemented")
}
