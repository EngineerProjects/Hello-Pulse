package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Provider implements the Provider interface for Amazon S3
type S3Provider struct {
	client *s3.Client
	config Config
}

// newS3Provider creates a new S3 provider
func newS3Provider(config Config) (*S3Provider, error) {
	return &S3Provider{
		config: config,
	}, nil
}

// Initialize initializes the S3 provider
func (p *S3Provider) Initialize(ctx context.Context) error {
	// Create AWS credentials using access and secret keys
	credProvider := credentials.NewStaticCredentialsProvider(p.config.AccessKey, p.config.SecretKey, "")

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credProvider),
		config.WithRegion(p.config.Region),
	)
	if err != nil {
		return fmt.Errorf("failed to load S3 configuration: %w", err)
	}

	// Create S3 client
	p.client = s3.NewFromConfig(cfg)

	// Test connection by listing buckets
	_, err = p.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("failed to connect to S3: %w", err)
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

// UploadFile uploads a file to S3
func (p *S3Provider) UploadFile(ctx context.Context, bucket string, path string, reader io.Reader, size int64, contentType string) (string, error) {
	// Create uploader
	uploader := manager.NewUploader(p.client)

	// Upload the file
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(path),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return path, nil
}

// DownloadFile downloads a file from S3
func (p *S3Provider) DownloadFile(ctx context.Context, bucket string, path string) (io.ReadCloser, error) {
	// Get object
	result, err := p.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}

	return result.Body, nil
}

// DeleteFile deletes a file from S3
func (p *S3Provider) DeleteFile(ctx context.Context, bucket string, path string) error {
	// Delete object
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// GetFileURL generates a presigned URL for accessing a file
func (p *S3Provider) GetFileURL(ctx context.Context, bucket string, path string, expires time.Duration) (string, error) {
	// Create presigner
	presignClient := s3.NewPresignClient(p.client)

	// Create presigned URL
	presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %w", err)
	}

	return presignedReq.URL, nil
}

// CreateBucket creates a new bucket
func (p *S3Provider) CreateBucket(ctx context.Context, bucket string) error {
	// Create bucket
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	// Add location constraint if the region is not us-east-1
	if p.config.Region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(p.config.Region),
		}
	}

	_, err := p.client.CreateBucket(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	return nil
}

// BucketExists checks if a bucket exists
func (p *S3Provider) BucketExists(ctx context.Context, bucket string) (bool, error) {
	// Check if bucket exists
	_, err := p.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		// If the bucket doesn't exist, AWS returns an error
		return false, nil
	}

	return true, nil
}