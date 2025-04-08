package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// Connect establishes a connection to the database
func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Paris",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ensureExtensions(db)

	DB = db
	log.Println("Successfully connected to database")
}

// ensureExtensions creates necessary PostgreSQL extensions
func ensureExtensions(db *gorm.DB) {
	// Create UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Fatalf("Failed to create uuid-ossp extension: %v", err)
	}
	log.Println("uuid-ossp extension ensured")

	// Create vector extension for AI features
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector;").Error; err != nil {
		log.Println("Warning: pgvector extension failed to load. Make sure it's installed if you need vector capabilities.")
	} else {
		log.Println("pgvector extension ensured")
	}
}

// RunMigrations runs the database migrations for all models
func RunMigrations(models ...interface{}) {
	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("Error during migration: %v", err)
	}
	log.Println("All tables migrated successfully")
}