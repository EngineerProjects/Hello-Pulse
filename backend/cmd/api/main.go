package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"hello-pulse.fr/internal/api/routes"
	"hello-pulse.fr/internal/models/auth"
	"hello-pulse.fr/internal/models/event"
	"hello-pulse.fr/internal/models/file"
	"hello-pulse.fr/internal/models/invite"
	"hello-pulse.fr/internal/models/organization"
	"hello-pulse.fr/internal/models/project"
	"hello-pulse.fr/internal/models/user"
	authrepo "hello-pulse.fr/internal/repositories/auth"
	eventrepo "hello-pulse.fr/internal/repositories/event"
	inviterepo "hello-pulse.fr/internal/repositories/invite"
	orgrepo "hello-pulse.fr/internal/repositories/organization"
	projectrepo "hello-pulse.fr/internal/repositories/project"
	userrepo "hello-pulse.fr/internal/repositories/user"
	authservice "hello-pulse.fr/internal/services/auth"
	eventservice "hello-pulse.fr/internal/services/event"
	orgservice "hello-pulse.fr/internal/services/organization"
	projectservice "hello-pulse.fr/internal/services/project"
	"hello-pulse.fr/pkg/config"
	"hello-pulse.fr/pkg/database"
)

func main() {
	// Load configuration
	appConfig := config.LoadConfig()

	// Connect to database
	database.Connect()

	// Run migrations
	database.RunMigrations(
		&user.User{},
		&organization.Organization{},
		&project.Project{},
		&auth.Session{},
		&event.Event{},
		&file.File{},
		&invite.InviteCode{},
	)

	// Initialize repositories
	userRepository := userrepo.NewRepository(database.DB)
	sessionRepository := authrepo.NewRepository(database.DB)
	orgRepository := orgrepo.NewRepository(database.DB)
	inviteRepository := inviterepo.NewRepository(database.DB)
	projectRepository := projectrepo.NewRepository(database.DB)
	eventRepository := eventrepo.NewRepository(database.DB)

	// Initialize services
	authService := authservice.NewService(userRepository, sessionRepository)
	projectService := projectservice.NewService(projectRepository, userRepository)
	orgService := orgservice.NewService(orgRepository, userRepository, inviteRepository)
	eventService := eventservice.NewService(eventRepository, userRepository)

	// Initialize Gin router
	r := gin.Default()
	r.MaxMultipartMemory = 30 << 30 // 30 GiB

	// Setup routes
	routes.Setup(r, authService, projectService, orgService, eventService)

	// Start server
	log.Printf("Starting server on port %s", appConfig.Port)
	if err := r.Run(":" + appConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}