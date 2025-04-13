package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient is a wrapper around minio.Client
type MinioClient struct {
	client *minio.Client
}

// MinioConfig holds the configuration for MinIO
type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

// NewMinioClient creates a new MinIO client
func NewMinioClient(config MinioConfig) (*MinioClient, error) {
	// Initialize the MinIO client
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Verify the connection by listing buckets
	_, err = client.ListBuckets(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	// Create default bucket if it doesn't exist
	err = ensureBucketExists(client, config.BucketName, "")
	if err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return &MinioClient{client: client}, nil
}

// ensureBucketExists checks if a bucket exists and creates it if it doesn't
func ensureBucketExists(client *minio.Client, bucketName, location string) error {
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
			Region: location,
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Bucket %s created successfully", bucketName)
	} else {
		log.Printf("Bucket %s already exists", bucketName)
	}

	return nil
}

// UploadFile uploads a file to MinIO
func (m *MinioClient) UploadFile(
	ctx context.Context,
	bucketName string,
	objectName string,
	reader io.Reader,
	fileSize int64,
	contentType string,
) (minio.UploadInfo, error) {
	// Set upload options
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	// Upload the file
	info, err := m.client.PutObject(ctx, bucketName, objectName, reader, fileSize, opts)
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("failed to upload file: %w", err)
	}

	return info, nil
}

// GetFileURL generates a presigned URL for a file
func (m *MinioClient) GetFileURL(
	ctx context.Context,
	bucketName string,
	objectName string,
	expires time.Duration,
) (string, error) {
	// Generate a presigned URL
	url, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expires, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// DeleteFile removes a file from MinIO
func (m *MinioClient) DeleteFile(ctx context.Context, bucketName string, objectName string) error {
	// Remove the file
	err := m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetFolderForFileType determines the appropriate folder based on file extension
func GetFolderForFileType(fileName string) string {
	ext := filepath.Ext(fileName)
	if ext == "" {
		return "others"
	}

	extLower := filepath.Ext(filepath.ToSlash(fileName))
	
	// Map of extensions to folder names
	folderMap := map[string]string{
		// Documents
		".pdf":  "documents",
		".doc":  "documents",
		".docx": "documents",
		".txt":  "documents",
		".rtf":  "documents",
		".odt":  "documents",
		".md":   "documents",
		".csv":  "documents",
		".xls":  "documents",
		".xlsx": "documents",
		".ppt":  "documents",
		".pptx": "documents",
		
		// Images
		".jpg":  "images",
		".jpeg": "images",
		".png":  "images",
		".gif":  "images",
		".webp": "images",
		".svg":  "images",
		".bmp":  "images",
		
		// Audio
		".mp3":  "audio",
		".wav":  "audio",
		".ogg":  "audio",
		".flac": "audio",
		
		// Video
		".mp4":  "video",
		".mov":  "video",
		".avi":  "video",
		".mkv":  "video",
		".webm": "video",
	}
	
	if folder, ok := folderMap[extLower]; ok {
		return folder
	}
	
	return "others"
}

// GenerateObjectName generates a unique object name for storage
func GenerateObjectName(organizationID uuid.UUID, folder, fileName string) string {
	// Generate a unique ID for the file to prevent name collisions
	uniqueID := uuid.New().String()
	
	// Create a path that includes organization, folder, and a unique ID
	extension := filepath.Ext(fileName)
	baseName := filepath.Base(fileName[:len(fileName)-len(extension)])
	
	return fmt.Sprintf("%s/%s/%s-%s%s", 
		organizationID.String(), 
		folder,
		baseName,
		uniqueID,
		extension,
	)
}