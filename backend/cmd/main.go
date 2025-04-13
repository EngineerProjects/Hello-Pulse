package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"hello-pulse.fr/internal/api/routes"
	"hello-pulse.fr/internal/models/auth"
	"hello-pulse.fr/internal/models/event"
	fileModel "hello-pulse.fr/internal/models/file"
	"hello-pulse.fr/internal/models/invite"
	"hello-pulse.fr/internal/models/organization"
	"hello-pulse.fr/internal/models/project"
	"hello-pulse.fr/internal/models/user"
	authrepo "hello-pulse.fr/internal/repositories/auth"
	eventrepo "hello-pulse.fr/internal/repositories/event"
	filerepo "hello-pulse.fr/internal/repositories/file"
	inviterepo "hello-pulse.fr/internal/repositories/invite"
	orgrepo "hello-pulse.fr/internal/repositories/organization"
	projectrepo "hello-pulse.fr/internal/repositories/project"
	userrepo "hello-pulse.fr/internal/repositories/user"
	authservice "hello-pulse.fr/internal/services/auth"
	eventservice "hello-pulse.fr/internal/services/event"
	fileservice "hello-pulse.fr/internal/services/file"
	orgservice "hello-pulse.fr/internal/services/organization"
	projectservice "hello-pulse.fr/internal/services/project"
	"hello-pulse.fr/pkg/config"
	"hello-pulse.fr/pkg/database"
	"hello-pulse.fr/pkg/storage"
)


func main() {
	// Load configuration
	appConfig := config.LoadConfig()
	storageConfig := config.LoadStorageConfig()

	// Connect to database
	database.Connect()

	// Run migrations
	database.RunMigrations(
		&user.User{},
		&organization.Organization{},
		&project.Project{},
		&project.Summary{},
		&auth.Session{},
		&event.Event{},
		&fileModel.File{},
		&invite.InviteCode{},
	)

	// Initialize MinIO client
	minioClient, err := storage.NewMinioClient(storage.MinioConfig{
		Endpoint:        storageConfig.Endpoint,
		AccessKeyID:     storageConfig.AccessKeyID,
		SecretAccessKey: storageConfig.SecretAccessKey,
		UseSSL:          storageConfig.UseSSL,
		BucketName:      "hello-pulse", // Default bucket name
	})
	
	if err != nil {
		log.Printf("Warning: Failed to initialize MinIO client: %v", err)
		log.Println("File storage functionality will be unavailable")
	}

	// Initialize repositories
	userRepository := userrepo.NewRepository(database.DB)
	sessionRepository := authrepo.NewRepository(database.DB)
	orgRepository := orgrepo.NewRepository(database.DB)
	inviteRepository := inviterepo.NewRepository(database.DB)
	projectRepository := projectrepo.NewRepository(database.DB)
	summaryRepository := projectrepo.NewSummaryRepository(database.DB)
	eventRepository := eventrepo.NewRepository(database.DB)
	fileRepository := filerepo.NewRepository(database.DB)

	// Initialize services
	authService := authservice.NewService(userRepository, sessionRepository)
	projectService := projectservice.NewService(projectRepository, userRepository, summaryRepository)
	orgService := orgservice.NewService(orgRepository, userRepository, inviteRepository)
	eventService := eventservice.NewService(eventRepository, userRepository)
	
	// Initialize file service only if MinIO client was successfully created
	var fileService *fileservice.Service
	if minioClient != nil {
		fileService = fileservice.NewService(fileRepository, minioClient, "hello-pulse")
		// Start the file cleanup background task
		StartFileCleanupTask(fileService)
	}

	// Initialize Gin router
	r := gin.Default()
	r.MaxMultipartMemory = 100 << 20 // 100 MiB for file uploads

	// Setup routes
	routes.Setup(
		r,
		database.DB,
		authService,
		projectService,
		orgService,
		eventService,
		fileService,
	)

	// Start server
	log.Printf("Starting server on port %s", appConfig.Port)
	if err := r.Run(":" + appConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}