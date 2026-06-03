package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"pos-system/backend/internal/infrastructure/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
	publicURL  string
}

func NewMinioStorage(cfg config.Config) (*MinioStorage, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.MinioBucketName)
	if err != nil {
		return nil, fmt.Errorf("check bucket exists: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.MinioBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("make bucket: %w", err)
		}

		// Set bucket policy to public read so uploaded files are accessible by the browser/frontend
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": "*",
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, cfg.MinioBucketName)
		err = client.SetBucketPolicy(ctx, cfg.MinioBucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("set bucket policy: %w", err)
		}
	}

	return &MinioStorage{
		client:     client,
		bucketName: cfg.MinioBucketName,
		publicURL:  cfg.MinioPublicURL,
	}, nil
}

func (s *MinioStorage) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Clean up public URL
	publicURL := strings.TrimSuffix(s.publicURL, "/")
	
	// If publicURL is empty (e.g. backend acts as a direct proxy), we could use relative path
	if publicURL == "" {
		return fmt.Sprintf("/uploads/products/%s", objectName), nil
	}
	
	// Check if publicURL contains schema, if not prefix with http://
	if !strings.HasPrefix(publicURL, "http://") && !strings.HasPrefix(publicURL, "https://") {
		publicURL = "http://" + publicURL
	}
	
	return fmt.Sprintf("%s/%s/%s", publicURL, s.bucketName, objectName), nil
}
