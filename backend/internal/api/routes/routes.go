// internal/api/routes/routes.go
package routes

import (
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
	"hello-pulse.fr/pkg/security"
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
	securityService *security.AuthorizationService,
) {
	// Create handlers
	authHandler := auth.NewHandler(authService)
	projectHandler := project.NewHandler(projectService, securityService)
	orgHandler := organization.NewHandler(orgService, securityService)
	eventHandler := event.NewHandler(eventService, securityService)
	fileHandler := file.NewHandler(fileService, securityService)

	// Create summary handler - add this line
	summaryRepository := projectrepo.NewSummaryRepository(db)
	summaryService := projectservice.NewSummaryService(summaryRepository, projectrepo.NewRepository(db))
	summaryHandler := project.NewSummaryHandler(summaryService, projectService, securityService)


	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.POST("/logout", authHandler.Logout)

	// Authentication middleware
	authMiddleware := middleware.AuthMiddleware(authService)
	
	// Security middleware - adds security context to requests
	securityMiddleware := middleware.SecurityMiddleware(securityService)
	
	// Organization required middleware
	orgRequiredMiddleware := middleware.OrganizationRequiredMiddleware(securityService)
	
	// Admin required middleware
	adminRequiredMiddleware := middleware.AdminRequiredMiddleware(securityService)

	// Basic authentication
	protected := router.Group("/", authMiddleware, securityMiddleware)
	
	// Organization-scoped routes
	orgProtected := protected.Group("/", orgRequiredMiddleware)

	// Configure routes
	{
		// User routes
		protected.GET("/me", authHandler.Me)

		// Organization routes
		protected.POST("/organizations", orgHandler.CreateOrganization)
		protected.POST("/organizations/join", orgHandler.JoinOrganization)
		
		// Organization routes that require organization membership
		orgProtected.GET("/organizations/invite-codes", adminRequiredMiddleware, orgHandler.GetInviteCodes)
		orgProtected.POST("/organizations/invite-codes", adminRequiredMiddleware, orgHandler.CreateInviteCode)
		orgProtected.DELETE("/organizations/invite-codes", adminRequiredMiddleware, orgHandler.DeleteInviteCode)

		// Project routes
		orgProtected.POST("/projects", projectHandler.CreateProject)
		orgProtected.GET("/projects", projectHandler.GetProjects)
		orgProtected.GET("/projects/:id", projectHandler.GetProject)
		orgProtected.PUT("/projects/:id", middleware.ResourceOwnerMiddleware(securityService, "project"), projectHandler.UpdateProject)
		orgProtected.DELETE("/projects/:id", middleware.ResourceOwnerMiddleware(securityService, "project"), projectHandler.DeleteProject)
		orgProtected.POST("/projects/add-user", projectHandler.AddParticipant)

		// Summary routes
		orgProtected.POST("/projects/summaries", summaryHandler.CreateSummary)
		orgProtected.GET("/projects/:id/summaries", summaryHandler.GetProjectSummaries)
		orgProtected.GET("/projects/summaries/:id", summaryHandler.GetSummary)
		orgProtected.PUT("/projects/summaries/:id", summaryHandler.UpdateSummary)
		orgProtected.DELETE("/projects/summaries/:id", summaryHandler.DeleteSummary)
		
		// Event routes
		orgProtected.POST("/events", eventHandler.CreateEvent)
		orgProtected.GET("/events", eventHandler.GetEvents)
		orgProtected.DELETE("/events/:id", middleware.ResourceOwnerMiddleware(securityService, "event"), eventHandler.DeleteEvent)
		orgProtected.POST("/events/add-member", eventHandler.AddParticipant)
		orgProtected.POST("/events/remove-member", eventHandler.RemoveParticipant)
		orgProtected.POST("/events/:id/update-title", middleware.ResourceOwnerMiddleware(securityService, "event"), eventHandler.UpdateEventTitle)
		orgProtected.GET("/events/:id/participants", eventHandler.GetEventParticipants)
		orgProtected.GET("/events/fetch-users", eventHandler.GetOrganizationUsers)
		
		// File routes
		orgProtected.POST("/files", fileHandler.UploadFile)
		orgProtected.GET("/files", fileHandler.GetUserFiles)
		orgProtected.GET("/files/organization", fileHandler.GetOrganizationFiles)
		orgProtected.GET("/files/types", fileHandler.GetFileTypes)
		orgProtected.GET("/files/:id", fileHandler.GetFileByID)
		orgProtected.GET("/files/:id/url", fileHandler.GetFileURL)
		orgProtected.DELETE("/files/:id", fileHandler.SoftDeleteFile)
		orgProtected.POST("/files/:id/restore", fileHandler.RestoreFile)
		orgProtected.POST("/files/batch-delete", fileHandler.BatchDeleteFiles)
		orgProtected.POST("/files/cleanup", adminRequiredMiddleware, fileHandler.RunCleanup)
		orgProtected.PUT("/files/:id/visibility", fileHandler.UpdateFileVisibility)
		orgProtected.GET("/files/:id/download", fileHandler.DownloadFile)
	}
}