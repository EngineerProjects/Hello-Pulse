package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioProvider implements the Provider interface for MinIO
type MinioProvider struct {
	client *minio.Client
	config Config
}

// newMinioProvider creates a new MinIO provider
func newMinioProvider(config Config) (*MinioProvider, error) {
	return &MinioProvider{
		config: config,
	}, nil
}

// Initialize initializes the MinIO provider
func (p *MinioProvider) Initialize(ctx context.Context) error {
	// Initialize MinIO client
	client, err := minio.New(p.config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(p.config.AccessKey, p.config.SecretKey, ""),
		Secure: p.config.UseSSL,
		Region: p.config.Region,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %w", err)
	}

	p.client = client

	// Test connection by listing buckets
	_, err = client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to MinIO server: %w", err)
	}

	// Create default bucket if specified
	if p.config.DefaultBucket != "" {
		exists, err := p.BucketExists(ctx, p.config.DefaultBucket)
		if err != nil {
			return fmt.Errorf("failed to check if bucket exists: %w", err)
		}

		if !exists {
			err = p.CreateBucket(ctx, p.config.DefaultBucket)
			if err != nil {
				return fmt.Errorf("failed to create default bucket: %w", err)
			}
		}
	}

	return nil
}

// UploadFile uploads a file to MinIO
func (p *MinioProvider) UploadFile(ctx context.Context, bucket string, path string, reader io.Reader, size int64, contentType string) (string, error) {
	options := minio.PutObjectOptions{
		ContentType: contentType,
	}

	// Upload the file
	_, err := p.client.PutObject(ctx, bucket, path, reader, size, options)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	return path, nil
}

// DownloadFile downloads a file from MinIO
func (p *MinioProvider) DownloadFile(ctx context.Context, bucket string, path string) (io.ReadCloser, error) {
	object, err := p.client.GetObject(ctx, bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from MinIO: %w", err)
	}

	return object, nil
}

// DeleteFile deletes a file from MinIO
func (p *MinioProvider) DeleteFile(ctx context.Context, bucket string, path string) error {
	err := p.client.RemoveObject(ctx, bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	return nil
}

// GetFileURL generates a URL for accessing a file
func (p *MinioProvider) GetFileURL(ctx context.Context, bucket string, path string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := p.client.PresignedGetObject(ctx, bucket, path, expires, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

// CreateBucket creates a new bucket
func (p *MinioProvider) CreateBucket(ctx context.Context, bucket string) error {
	err := p.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
		Region: p.config.Region,
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	return nil
}

// BucketExists checks if a bucket exists
func (p *MinioProvider) BucketExists(ctx context.Context, bucket string) (bool, error) {
	exists, err := p.client.BucketExists(ctx, bucket)
	if err != nil {
		return false, fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	return exists, nil
}