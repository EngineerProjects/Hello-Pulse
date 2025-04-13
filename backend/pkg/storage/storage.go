package storage

import (
	"context"
	"io"
	"time"
)

// Provider defines the interface that all storage providers must implement
type Provider interface {
	// Initialize initializes the storage provider
	Initialize(ctx context.Context) error
	// UploadFile uploads a file to storage
	UploadFile(ctx context.Context, bucket string, path string, reader io.Reader, size int64, contentType string) (string, error)
	// DownloadFile downloads a file from storage
	DownloadFile(ctx context.Context, bucket string, path string) (io.ReadCloser, error)
	// DeleteFile deletes a file from storage
	DeleteFile(ctx context.Context, bucket string, path string) error
	// GetFileURL generates a URL for accessing a file
	GetFileURL(ctx context.Context, bucket string, path string, expires time.Duration) (string, error)
	// CreateBucket creates a new bucket
	CreateBucket(ctx context.Context, bucket string) error
	// BucketExists checks if a bucket exists
	BucketExists(ctx context.Context, bucket string) (bool, error)
}

// Config holds the storage configuration
type Config struct {
	Provider      string            // e.g., "minio", "s3", "azure"
	Endpoint      string
	Region        string
	UseSSL        bool
	AccessKey     string
	SecretKey     string
	DefaultBucket string
	Options       map[string]string // Additional provider-specific options
}

// FileCategory represents the category of a file
type FileCategory string

const (
	CategoryDocument FileCategory = "documents"
	CategoryImage    FileCategory = "images"
	CategoryAudio    FileCategory = "audio"
	CategoryVideo    FileCategory = "video"
	CategoryArchive  FileCategory = "archives"
	CategoryOther    FileCategory = "others"
)