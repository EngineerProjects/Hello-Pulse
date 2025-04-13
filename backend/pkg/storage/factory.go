// pkg/storage/factory.go
package storage

import (
	"fmt"
)

// NewProvider creates a storage provider based on the provider specified in the config
func NewProvider(config Config) (Provider, error) {
	switch config.Provider {
	case "minio":
		return newMinioProvider(config)
	// Add more cases here for other storage providers
	case "s3":
	    return newS3Provider(config)
	// case "azure":
	//     return newAzureProvider(config)
	// case "gcs":
	//     return newGCSProvider(config)
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", config.Provider)
	}
}