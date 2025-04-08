package event

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/event"
	"hello-pulse.fr/internal/models/user"
	eventrepo "hello-pulse.fr/internal/repositories/event"
	userrepo "hello-pulse.fr/internal/repositories/user"
)

// Service handles event business logic
type Service struct {
	eventRepo *eventrepo.Repository
	userRepo  *userrepo.Repository
}

// NewService creates a new event service
func NewService(eventRepo *eventrepo.Repository, userRepo *userrepo.Repository) *Service {
	return &Service{
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

// CreateEvent creates a new event
func (s *Service) CreateEvent(
	title string,
	date time.Time,
	startTime time.Time,
	endTime time.Time,
	importance string,
	creatorID uuid.UUID,
	orgID uuid.UUID,
	userIDs []uuid.UUID,
) (*event.Event, error) {
	// Validate creator
	creator, err := s.userRepo.FindByID(creatorID)
	if err != nil {
		return nil, errors.New("creator not found")
	}

	if creator.OrganizationID == nil || *creator.OrganizationID != orgID {
		return nil, errors.New("creator does not belong to the specified organization")
	}

	// Validate that all users belong to the same organization
	for _, userID := range userIDs {
		if userID == creatorID {
			return nil, errors.New("creator cannot be added as a participant")
		}

		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return nil, errors.New("user not found")
		}

		if user.OrganizationID == nil || *user.OrganizationID != orgID {
			return nil, errors.New("one or more users do not belong to the same organization")
		}
	}

	// Create event
	newEvent := &event.Event{
		Title:          title,
		Date:           date,
		StartTime:      startTime,
		EndTime:        endTime,
		OrganizationID: orgID,
		CreatedByID:    creatorID,
		Importance:     importance,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.eventRepo.Create(newEvent); err != nil {
		return nil, err
	}

	// Add participants
	for _, userID := range userIDs {
		if err := s.eventRepo.AddParticipant(newEvent.EventID, userID); err != nil {
			return nil, err
		}
	}

	return newEvent, nil
}

// GetEvent retrieves an event by ID
func (s *Service) GetEvent(eventID uuid.UUID) (*event.Event, error) {
	return s.eventRepo.FindByID(eventID)
}

// GetOrganizationEvents retrieves all events for an organization
func (s *Service) GetOrganizationEvents(orgID uuid.UUID) ([]event.Event, error) {
	return s.eventRepo.FindByOrganization(orgID)
}

// GetUserEvents retrieves all events for a user
func (s *Service) GetUserEvents(userID uuid.UUID) ([]event.Event, error) {
	return s.eventRepo.FindForUser(userID)
}

// UpdateEventTitle updates an event's title
func (s *Service) UpdateEventTitle(eventID uuid.UUID, title string) error {
	return s.eventRepo.UpdateTitle(eventID, title)
}

// DeleteEvent deletes an event
func (s *Service) DeleteEvent(eventID uuid.UUID) error {
	// Clear participants
	if err := s.eventRepo.ClearParticipants(eventID); err != nil {
		return err
	}

	// Delete the event
	return s.eventRepo.Delete(eventID)
}

// AddParticipant adds a user to an event
func (s *Service) AddParticipant(eventID, userID uuid.UUID) error {
	// Verify event exists
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if user belongs to the event's organization
	if user.OrganizationID == nil || *user.OrganizationID != event.OrganizationID {
		return errors.New("user does not belong to the event's organization")
	}

	// Check if user is not the creator (creator is not a participant)
	if event.CreatedByID == userID {
		return errors.New("cannot add event creator as a participant")
	}

	// Add participant
	return s.eventRepo.AddParticipant(eventID, userID)
}

// RemoveParticipant removes a user from an event
func (s *Service) RemoveParticipant(eventID, userID uuid.UUID) error {
	return s.eventRepo.RemoveParticipant(eventID, userID)
}

// GetEventParticipants retrieves all participants of an event
func (s *Service) GetEventParticipants(eventID uuid.UUID) ([]user.User, error) {
	return s.eventRepo.GetParticipants(eventID)
}

// GetUpcomingEvents retrieves upcoming events
func (s *Service) GetUpcomingEvents(orgID uuid.UUID) ([]event.Event, error) {
	today := time.Now().Truncate(24 * time.Hour) // Start of today
	return s.eventRepo.FindUpcoming(orgID, today)
}