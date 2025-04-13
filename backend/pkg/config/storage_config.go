package config

import (
	"strconv"

	"github.com/joho/godotenv"
	"hello-pulse.fr/pkg/storage"
)

// LoadStorageConfig loads the storage configuration
func LoadStorageConfig() storage.Config {
	_ = godotenv.Load() // Ignoring error as it's handled elsewhere
	
	useSSL, _ := strconv.ParseBool(GetEnv("MINIO_USE_SSL", "false"))
	
	return storage.Config{
		Provider:      GetEnv("STORAGE_PROVIDER", "minio"), // Default provider
		Endpoint:      GetEnv("MINIO_ENDPOINT", "minio:9000"),
		Region:        GetEnv("STORAGE_REGION", ""),
		UseSSL:        useSSL,
		AccessKey:     GetEnv("MINIO_ROOT_USER", ""),
		SecretKey:     GetEnv("MINIO_ROOT_PASSWORD", ""),
		DefaultBucket: GetEnv("STORAGE_DEFAULT_BUCKET", "hello-pulse"),
		Options:       map[string]string{},
	}
}