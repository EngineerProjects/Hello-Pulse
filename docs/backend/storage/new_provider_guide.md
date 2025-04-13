# Adding New Storage Providers Guide

This guide explains how to add support for additional storage providers to the Hello Pulse platform.

## Architecture Overview

The storage abstraction is designed to be provider-agnostic, allowing you to easily switch between different providers:

- **Interface**: The `storage.Provider` interface defines the operations a storage provider must implement.
- **Factory**: The `storage.NewProvider` function creates the appropriate provider based on configuration.
- **Implementations**: Individual provider implementations (like `MinioProvider`, `S3Provider`, etc.) 

## Adding a New Provider

To add a new storage provider, follow these steps:

### 1. Create a new file for your provider

Create a new file in the `pkg/storage` directory, e.g., `azure_provider.go`.

### 2. Implement the Provider interface

Your implementation must fulfill all methods in the `storage.Provider` interface:

```go
// pkg/storage/azure_provider.go
package storage

import (
	"context"
	"fmt"
	"io"
	"time"
	
	// Import Azure SDK packages
)

// AzureProvider implements the Provider interface for Azure Blob Storage
type AzureProvider struct {
	// Add Azure client here
	config Config
}

// newAzureProvider creates a new Azure Blob Storage provider
func newAzureProvider(config Config) (*AzureProvider, error) {
	return &AzureProvider{
		config: config,
	}, nil
}

// Implement all methods required by the Provider interface
func (p *AzureProvider) Initialize(ctx context.Context) error {
	// Initialize Azure client
}

func (p *AzureProvider) UploadFile(ctx context.Context, bucket string, path string, reader io.Reader, size int64, contentType string) (string, error) {
	// Upload file to Azure Blob Storage
}

// ... Implement all other required methods
```

### 3. Update the factory function

Add your new provider to the factory function in `pkg/storage/factory.go`:

```go
// NewProvider creates a storage provider based on the provider specified in the config
func NewProvider(config Config) (Provider, error) {
	switch config.Provider {
	case "minio":
		return newMinioProvider(config)
	case "s3":
		return newS3Provider(config)
	case "azure":
		return newAzureProvider(config)
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", config.Provider)
	}
}
```

### 4. Add provider-specific configuration options

For provider-specific settings, use the `Options` map in the config:

```go
// In your .env file
STORAGE_PROVIDER=azure
AZURE_ACCOUNT_NAME=mystorageaccount
AZURE_ACCOUNT_KEY=myaccountkey

// In your provider implementation
func (p *AzureProvider) Initialize(ctx context.Context) error {
	accountName, ok := p.config.Options["AccountName"]
	if !ok {
		accountName = p.config.AccessKey // Fall back to AccessKey if not provided
	}
	
	accountKey, ok := p.config.Options["AccountKey"]
	if !ok {
		accountKey = p.config.SecretKey // Fall back to SecretKey if not provided
	}
	
	// Use these values to initialize your client
}
```

### 5. Update the configuration loading

Modify `pkg/config/storage_config.go` to support your new provider:

```go
// LoadStorageConfig loads the storage configuration
func LoadStorageConfig() storage.Config {
	// ... existing code ...
	
	// Add provider-specific options
	options := map[string]string{}
	
	if GetEnv("STORAGE_PROVIDER", "minio") == "azure" {
		options["AccountName"] = GetEnv("AZURE_ACCOUNT_NAME", "")
		options["AccountKey"] = GetEnv("AZURE_ACCOUNT_KEY", "")
	}
	
	return storage.Config{
		Provider:      GetEnv("STORAGE_PROVIDER", "minio"),
		// ... existing fields ...
		Options:       options,
	}
}
```

## Provider Implementation Guidelines

When implementing a storage provider, follow these guidelines:

1. **Handle Bucket/Container Naming**: Different providers may have different naming conventions and restrictions for buckets or containers.

2. **URL Generation**: URL generation (presigned URLs) might work differently across providers.

3. **Error Handling**: Translate provider-specific errors into meaningful standard errors.

4. **Dependencies**: Declare dependencies in `go.mod` and import them in your implementation.

5. **Idempotency**: Operations should be idempotent where possible (e.g., creating a bucket that already exists should not error).

## Testing Your Implementation

Create tests for your new provider:

1. Create a file `pkg/storage/your_provider_test.go`
2. Test all provider methods against a real or mocked service
3. Include integration tests if possible

## Example: Google Cloud Storage Provider

```go
// pkg/storage/gcs_provider.go
package storage

import (
	"context"
	"fmt"
	"io"
	"time"
	
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// GCSProvider implements the Provider interface for Google Cloud Storage
type GCSProvider struct {
	client *storage.Client
	config Config
}

// newGCSProvider creates a new Google Cloud Storage provider
func newGCSProvider(config Config) (*GCSProvider, error) {
	return &GCSProvider{
		config: config,
	}, nil
}

// Initialize initializes the GCS provider
func (p *GCSProvider) Initialize(ctx context.Context) error {
	// Read credentials from config
	credentialsFile, ok := p.config.Options["CredentialsFile"]
	
	var client *storage.Client
	var err error
	
	if ok && credentialsFile != "" {
		// Initialize with credentials file
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	} else {
		// Use default credentials
		client, err = storage.NewClient(ctx)
	}
	
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %w", err)
	}
	
	p.client = client
	return nil
}

// Implement remaining methods
```

Remember to update the factory function and configuration loader to support your new provider.