package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	
	"hello-pulse.fr/internal/api/handlers/auth"
	"hello-pulse.fr/internal/api/handlers/event"
	"hello-pulse.fr/internal/api/handlers/file"
	"hello-pulse.fr/internal/api/handlers/organization"
	"hello-pulse.fr/internal/api/handlers/project"
	"hello-pulse.fr/internal/api/middleware"
	
	authservice "hello-pulse.fr/internal/services/auth"
	eventservice "hello-pulse.fr/internal/services/event"
	fileservice "hello-pulse.fr/internal/services/file"
	orgservice "hello-pulse.fr/internal/services/organization"
	projectservice "hello-pulse.fr/internal/services/project"
	
	projectrepo "hello-pulse.fr/internal/repositories/project"
)

// Setup configures all API routes
func Setup(
	router *gin.Engine,
	db *gorm.DB,
	authService *authservice.Service,
	projectService *projectservice.Service,
	orgService *orgservice.Service,
	eventService *eventservice.Service,
	fileService *fileservice.Service,
) {
	// Create handlers
	authHandler := auth.NewHandler(authService)
	projectHandler := project.NewHandler(projectService)
	orgHandler := organization.NewHandler(orgService)
	eventHandler := event.NewHandler(eventService)

	// Create summary repository, service and handler
	summaryRepo := projectrepo.NewSummaryRepository(db)
	summaryService := projectservice.NewSummaryService(summaryRepo, projectrepo.NewRepository(db))
	summaryHandler := project.NewSummaryHandler(summaryService, projectService)

	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.POST("/logout", authHandler.Logout)

	// Protected routes
	authorized := router.Group("/", middleware.AuthMiddleware(authService))
	{
		// User routes
		authorized.GET("/me", authHandler.Me)

		// Organization routes
		authorized.POST("/organizations", orgHandler.CreateOrganization)
		authorized.POST("/organizations/join", orgHandler.JoinOrganization)
		authorized.GET("/organizations/invite-codes", orgHandler.GetInviteCodes)
		authorized.POST("/organizations/invite-codes", orgHandler.CreateInviteCode)
		authorized.DELETE("/organizations/invite-codes", orgHandler.DeleteInviteCode)

		// Project routes
		authorized.POST("/projects", projectHandler.CreateProject)
		authorized.GET("/projects", projectHandler.GetProjects)
		authorized.GET("/projects/:id", projectHandler.GetProject)
		authorized.PUT("/projects/:id", projectHandler.UpdateProject)
		authorized.DELETE("/projects/:id", projectHandler.DeleteProject)
		authorized.POST("/projects/add-user", projectHandler.AddParticipant)

		// Summary routes
		authorized.POST("/projects/summaries", summaryHandler.CreateSummary)
		authorized.GET("/projects/:id/summaries", summaryHandler.GetProjectSummaries)
		authorized.GET("/projects/summaries/:id", summaryHandler.GetSummary)
		authorized.PUT("/projects/summaries/:id", summaryHandler.UpdateSummary)
		authorized.DELETE("/projects/summaries/:id", summaryHandler.DeleteSummary)

		// Event routes
		authorized.POST("/events", eventHandler.CreateEvent)
		authorized.GET("/events", eventHandler.GetEvents)
		authorized.DELETE("/events/:id", eventHandler.DeleteEvent)
		authorized.POST("/events/add-member", eventHandler.AddParticipant)
		authorized.POST("/events/remove-member", eventHandler.RemoveParticipant)
		authorized.POST("/events/:id/update-title", eventHandler.UpdateEventTitle)
		authorized.GET("/events/:id/participants", eventHandler.GetEventParticipants)
		authorized.GET("/events/fetch-users", eventHandler.GetOrganizationUsers)
		
		// File routes - but only if the file service is available
		if fileService != nil {
			fileHandler := file.NewHandler(fileService)
			
			authorized.POST("/files", fileHandler.UploadFile)
			authorized.GET("/files", fileHandler.GetUserFiles)
			authorized.GET("/files/organization", fileHandler.GetOrganizationFiles)
			authorized.GET("/files/types", fileHandler.GetFileTypes)
			authorized.GET("/files/:id", fileHandler.GetFileByID)
			authorized.GET("/files/:id/url", fileHandler.GetFileURL)
			authorized.DELETE("/files/:id", fileHandler.SoftDeleteFile)
			authorized.POST("/files/:id/restore", fileHandler.RestoreFile)
			authorized.POST("/files/batch-delete", fileHandler.BatchDeleteFiles)
			authorized.POST("/files/cleanup", fileHandler.RunCleanup)
		} else {
			log.Println("File routes not configured: file service is unavailable")
		}
	}
}