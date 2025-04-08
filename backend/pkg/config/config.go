package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// AppConfig holds application configuration
type AppConfig struct {
	Port           string
	JWTSecret      string
	CleanupTimer   time.Duration
	MaxGoroutines  int
	FileExpiration time.Duration
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

// DBConfig holds database configuration
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DbName   string
}

// LoadConfig loads the application configuration
func LoadConfig() *AppConfig {
	_ = godotenv.Load() // Ignoring error as it's handled elsewhere

	maxGoroutines, _ := strconv.Atoi(GetEnv("MAX_GOROUTINES", "5"))
	
	return &AppConfig{
		Port:           GetEnv("PORT", "8000"),
		JWTSecret:      GetEnv("JWT_SECRET", "your-secret-key"),
		CleanupTimer:   10 * 24 * time.Hour, // 10 days default
		MaxGoroutines:  maxGoroutines,
		FileExpiration: 10 * time.Minute,
	}
}

// LoadStorageConfig loads the storage configuration
func LoadStorageConfig() *StorageConfig {
	useSSL, _ := strconv.ParseBool(GetEnv("MINIO_USE_SSL", "false"))
	
	return &StorageConfig{
		Endpoint:        GetEnv("MINIO_ENDPOINT", "minio:9000"),
		AccessKeyID:     GetEnv("MINIO_ROOT_USER", ""),
		SecretAccessKey: GetEnv("MINIO_ROOT_PASSWORD", ""),
		UseSSL:          useSSL,
	}
}

// LoadDBConfig loads the database configuration
func LoadDBConfig() *DBConfig {
	return &DBConfig{
		User:     GetEnv("POSTGRES_USER", ""),
		Password: GetEnv("POSTGRES_PASSWORD", ""),
		Host:     GetEnv("POSTGRES_HOST", "localhost"),
		Port:     GetEnv("POSTGRES_PORT", "5432"),
		DbName:   GetEnv("POSTGRES_DB_NAME", "hellopulsedb"),
	}
}

// GetEnv gets an environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}