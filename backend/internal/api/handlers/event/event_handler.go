package event

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/user"
	"hello-pulse.fr/internal/services/event"
	"hello-pulse.fr/pkg/security"
)

// Handler handles event API endpoints
type Handler struct {
	eventService *event.Service
	securityService *security.AuthorizationService 
}

// NewHandler creates a new event handler
func NewHandler(eventService *event.Service, securityService *security.AuthorizationService) *Handler {
	return &Handler{
		eventService: eventService,
		securityService: securityService,
	}
}

// CreateEventRequest represents the create event request payload
type CreateEventRequest struct {
	Title      string   `json:"title" binding:"required"`
	Date       string   `json:"date" binding:"required"`     // Format: YYYY-MM-DD
	StartTime  string   `json:"startTime" binding:"required"` // Format: HH:MM
	EndTime    string   `json:"endTime" binding:"required"`   // Format: HH:MM
	UserIDs    []string `json:"userIds"`                     // List of user IDs to add as participants
	Importance string   `json:"importance" binding:"required"` // Event importance level
}

// UpdateEventTitleRequest represents the update event title request payload
type UpdateEventTitleRequest struct {
	EventID string `json:"eventId" binding:"required"`
	Title   string `json:"title" binding:"required"`
}

// ParticipantRequest represents the add/remove participant request payload
type ParticipantRequest struct {
	EventID string `json:"eventId" binding:"required"`
	UserID  string `json:"userId" binding:"required"`
}

// CreateEvent handles event creation
func (h *Handler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	// Parse date and time
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid start time format. Use HH:MM",
		})
		return
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid end time format. Use HH:MM",
		})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	if user.OrganizationID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User does not belong to an organization",
		})
		return
	}

	// Parse user IDs
	var userIDs []uuid.UUID
	for _, idStr := range req.UserIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid user ID format",
			})
			return
		}
		userIDs = append(userIDs, id)
	}

	// Create event
	newEvent, err := h.eventService.CreateEvent(
		req.Title,
		date,
		startTime,
		endTime,
		req.Importance,
		user.UserID,
		*user.OrganizationID,
		userIDs,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event created successfully",
		"event": gin.H{
			"id":        newEvent.EventID,
			"title":     newEvent.Title,
			"date":      newEvent.Date.Format("2006-01-02"),
			"startTime": newEvent.StartTime.Format("15:04"),
			"endTime":   newEvent.EndTime.Format("15:04"),
			"importance": newEvent.Importance,
		},
	})
}

// GetEvents handles retrieving all events for the user
func (h *Handler) GetEvents(c *gin.Context) {
	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Get events
	events, err := h.eventService.GetUserEvents(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve events",
		})
		return
	}

	// Format events for response
	var formattedEvents []gin.H
	for _, e := range events {
		formattedEvents = append(formattedEvents, gin.H{
			"id":           e.EventID,
			"title":        e.Title,
			"date":         e.Date.Format("2006-01-02"),
			"startTime":    e.StartTime.Format("15:04"),
			"endTime":      e.EndTime.Format("15:04"),
			"importance":   e.Importance,
			"createdById":  e.CreatedByID,
			"organizationId": e.OrganizationID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"events":  formattedEvents,
		"userId":  user.UserID,
	})
}

// DeleteEvent handles deleting an event
func (h *Handler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	eventID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid event ID",
		})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user can modify this event
	canModify, err := h.securityService.CanModifyEvent(c.Request.Context(), user.UserID, eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Event not found",
		})
		return
	}
	
	if !canModify {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only the event creator can delete the event",
		})
		return
	}

	// Delete event
	if err := h.eventService.DeleteEvent(eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event deleted successfully",
	})
}

// AddParticipant handles adding a user to an event
func (h *Handler) AddParticipant(c *gin.Context) {
	var req ParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid event ID",
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user can modify this event
	canModify, err := h.securityService.CanModifyEvent(c.Request.Context(), user.UserID, eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Event not found",
		})
		return
	}
	
	if !canModify {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only the event creator can add participants",
		})
		return
	}

	// Add participant
	if err := h.eventService.AddParticipant(eventID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Participant added successfully",
	})
}


// RemoveParticipant handles removing a user from an event
func (h *Handler) RemoveParticipant(c *gin.Context) {
	var req ParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid event ID",
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user can modify this event
	canModify, err := h.securityService.CanModifyEvent(c.Request.Context(), user.UserID, eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Event not found",
		})
		return
	}
	
	if !canModify {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only the event creator can remove participants",
		})
		return
	}

	// Remove participant
	if err := h.eventService.RemoveParticipant(eventID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to remove participant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Participant removed successfully",
	})
}

// UpdateEventTitle handles updating an event's title
func (h *Handler) UpdateEventTitle(c *gin.Context) {
	var req UpdateEventTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid event ID",
		})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user can modify this event
	canModify, err := h.securityService.CanModifyEvent(c.Request.Context(), user.UserID, eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Event not found",
		})
		return
	}
	
	if !canModify {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only the event creator can update the event",
		})
		return
	}

	// Update event title
	if err := h.eventService.UpdateEventTitle(eventID, req.Title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update event title",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event title updated successfully",
	})
}

// GetEventParticipants handles retrieving participants of an event
func (h *Handler) GetEventParticipants(c *gin.Context) {
	id := c.Param("id")
	eventID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid event ID",
		})
		return
	}

	// Get event
	event, err := h.eventService.GetEvent(eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Event not found",
		})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user is part of the same organization
	if user.OrganizationID == nil || *user.OrganizationID != event.OrganizationID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Access denied",
		})
		return
	}

	// Get participants
	participants, err := h.eventService.GetEventParticipants(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve participants",
		})
		return
	}

	// Format participants for response
	var formattedParticipants []gin.H
	for _, p := range participants {
		formattedParticipants = append(formattedParticipants, gin.H{
			"id":        p.UserID,
			"firstName": p.FirstName,
			"lastName":  p.LastName,
			"email":     p.Email,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"participants": formattedParticipants,
	})
}

// GetOrganizationUsers handles retrieving users from the organization for event creation
func (h *Handler) GetOrganizationUsers(c *gin.Context) {
	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	if user.OrganizationID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User does not belong to an organization",
		})
		return
	}

	// This endpoint relies on a user service to get organization users
	// For now, we'll return a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "To implement: Get organization users",
	})
}