package routes

import (
	"github.com/gin-gonic/gin"
	"hello-pulse.fr/internal/api/handlers/auth"
	"hello-pulse.fr/internal/api/handlers/event"
	"hello-pulse.fr/internal/api/handlers/organization"
	"hello-pulse.fr/internal/api/handlers/project"
	"hello-pulse.fr/internal/api/middleware"
	authservice "hello-pulse.fr/internal/services/auth"
	eventservice "hello-pulse.fr/internal/services/event"
	orgservice "hello-pulse.fr/internal/services/organization"
	projectservice "hello-pulse.fr/internal/services/project"
)

// Setup configures all API routes
func Setup(
	router *gin.Engine,
	authService *authservice.Service,
	projectService *projectservice.Service,
	orgService *orgservice.Service,
	eventService *eventservice.Service,
) {
	// Create handlers
	authHandler := auth.NewHandler(authService)
	projectHandler := project.NewHandler(projectService)
	orgHandler := organization.NewHandler(orgService)
	eventHandler := event.NewHandler(eventService)

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

		// Project routes
		authorized.POST("/projects", projectHandler.CreateProject)
		authorized.GET("/projects", projectHandler.GetProjects)
		authorized.GET("/projects/:id", projectHandler.GetProject)
		authorized.PUT("/projects/:id", projectHandler.UpdateProject)
		authorized.DELETE("/projects/:id", projectHandler.DeleteProject)
		authorized.POST("/projects/add-user", projectHandler.AddParticipant)

		// Event routes
		authorized.POST("/events", eventHandler.CreateEvent)
		authorized.GET("/events", eventHandler.GetEvents)
		authorized.DELETE("/events/:id", eventHandler.DeleteEvent)
		authorized.POST("/events/add-member", eventHandler.AddParticipant)
		authorized.POST("/events/remove-member", eventHandler.RemoveParticipant)
		authorized.POST("/events/:id/update-title", eventHandler.UpdateEventTitle)
		authorized.GET("/events/:id/participants", eventHandler.GetEventParticipants)
		authorized.GET("/events/fetch-users", eventHandler.GetOrganizationUsers)
	}
}