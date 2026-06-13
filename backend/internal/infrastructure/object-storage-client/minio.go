package objectstorageclient

import (
	"context"
	"io"
	"strings"

	"pos-system/backend/internal/infrastructure/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client    *minio.Client
	bucket    string
	publicURL string
	useSSL    bool
}

func NewMinio(cfg config.Config) (*MinioStorage, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioStorage{
		client:    client,
		bucket:    cfg.MinioBucketName,
		publicURL: strings.TrimRight(cfg.MinioPublicURL, "/"),
		useSSL:    cfg.MinioUseSSL,
	}, nil
}

func (s *MinioStorage) UploadProductImage(ctx context.Context, filename string, reader io.Reader, size int64, contentType string) (string, error) {
	if err := s.ensureBucket(ctx); err != nil {
		return "", err
	}
	objectName := "products/" + filename
	if _, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return "", err
	}
	return s.objectURL(objectName), nil
}

func (s *MinioStorage) ensureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}
	return s.client.SetBucketPolicy(ctx, s.bucket, publicReadPolicy(s.bucket))
}

func (s *MinioStorage) objectURL(objectName string) string {
	if s.publicURL != "" {
		return s.publicURL + "/" + s.bucket + "/" + objectName
	}
	scheme := "http"
	if s.useSSL {
		scheme = "https"
	}
	return scheme + "://" + s.client.EndpointURL().Host + "/" + s.bucket + "/" + objectName
}

func publicReadPolicy(bucket string) string {
	return `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::` + bucket + `/products/*"]
    }
  ]
}`
}
