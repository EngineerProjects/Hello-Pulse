package security

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/user"
	fileRepo "hello-pulse.fr/internal/repositories/file"
	orgRepo "hello-pulse.fr/internal/repositories/organization"
	projectRepo "hello-pulse.fr/internal/repositories/project"
	eventRepo "hello-pulse.fr/internal/repositories/event"
	userRepo "hello-pulse.fr/internal/repositories/user"
)

var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")
	
	// ErrAccessDenied is returned when a user doesn't have access to a resource
	ErrAccessDenied = errors.New("access denied")
	
	// ErrNotMember is returned when a user is not a member of an organization
	ErrNotMember = errors.New("user is not a member of the organization")
	
	// ErrNotProjectParticipant is returned when a user is not a participant in a project
	ErrNotProjectParticipant = errors.New("user is not a participant in the project")
	
	// ErrNotEventParticipant is returned when a user is not a participant in an event
	ErrNotEventParticipant = errors.New("user is not a participant in the event")
	
	// ErrNotOwner is returned when a user is not the owner of a resource
	ErrNotOwner = errors.New("user is not the owner of this resource")
	
	// ErrNotAdmin is returned when a user is not an admin of the organization
	ErrNotAdmin = errors.New("user is not an admin of the organization")
)

// Role definitions
const (
	RoleAdmin = "Admin"
	RoleUser  = "User"
)

// AuthorizationService provides centralized security checks for the application
type AuthorizationService struct {
	fileRepo        *fileRepo.Repository
	projectRepo     *projectRepo.Repository
	orgRepo         *orgRepo.Repository
	userRepo        *userRepo.Repository
	eventRepo       *eventRepo.Repository
}

// NewAuthorizationService creates a new instance of the authorization service
func NewAuthorizationService(
	fileRepo *fileRepo.Repository,
	projectRepo *projectRepo.Repository,
	orgRepo *orgRepo.Repository,
	userRepo *userRepo.Repository,
	eventRepo *eventRepo.Repository,
) *AuthorizationService {
	return &AuthorizationService{
		fileRepo:    fileRepo,
		projectRepo: projectRepo,
		orgRepo:     orgRepo,
		userRepo:    userRepo,
		eventRepo:   eventRepo,
	}
}

// GetUser retrieves a user by ID
func (s *AuthorizationService) GetUser(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return user, nil
}

// IsUserInOrganization checks if a user belongs to an organization
func (s *AuthorizationService) IsUserInOrganization(ctx context.Context, userID, orgID uuid.UUID) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	return user.OrganizationID != nil && *user.OrganizationID == orgID, nil
}

// IsUserAdmin checks if a user is an admin of their organization
func (s *AuthorizationService) IsUserAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	return user.Role == RoleAdmin, nil
}

// CanAccessFile checks if a user can access a file
func (s *AuthorizationService) CanAccessFile(ctx context.Context, userID, fileID uuid.UUID) (bool, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve file: %w", err)
	}
	
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	// User can access if:
	// 1. They're the uploader
	// 2. The file is public AND they belong to the same organization
	// 3. They belong to the same organization (for non-public files)
	
	// Check if user is the uploader
	if file.UploaderID == userID {
		return true, nil
	}
	
	// Check organization membership
	isSameOrg := user.OrganizationID != nil && *user.OrganizationID == file.OrganizationID
	if !isSameOrg {
		return false, nil
	}
	
	// For non-public files, user must be in the same organization
	return file.IsPublic || isSameOrg, nil
}

// CanModifyFile checks if a user can modify a file (update, delete, restore)
func (s *AuthorizationService) CanModifyFile(ctx context.Context, userID, fileID uuid.UUID) (bool, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve file: %w", err)
	}
	
	// Only the uploader can modify the file
	return file.UploaderID == userID, nil
}

// CanAccessProject checks if a user can access a project
func (s *AuthorizationService) CanAccessProject(ctx context.Context, userID, projectID uuid.UUID) (bool, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve project: %w", err)
	}
	
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	// User can access if:
	// 1. User belongs to the same organization as the project
	if user.OrganizationID == nil || *user.OrganizationID != project.OrganizationID {
		return false, nil
	}
	
	// Check if user is the owner or a participant
	if project.OwnerID == userID {
		return true, nil
	}
	
	// Check if user is a participant
	participants, err := s.projectRepo.GetParticipants(projectID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve project participants: %w", err)
	}
	
	for _, participant := range participants {
		if participant.UserID == userID {
			return true, nil
		}
	}
	
	return false, nil
}

// CanModifyProject checks if a user can modify a project (update, delete)
func (s *AuthorizationService) CanModifyProject(ctx context.Context, userID, projectID uuid.UUID) (bool, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve project: %w", err)
	}
	
	// Only the owner can modify the project
	return project.OwnerID == userID, nil
}

// CanAccessEvent checks if a user can access an event
func (s *AuthorizationService) CanAccessEvent(ctx context.Context, userID, eventID uuid.UUID) (bool, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve event: %w", err)
	}
	
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	// User can access if:
	// 1. User belongs to the same organization as the event
	if user.OrganizationID == nil || *user.OrganizationID != event.OrganizationID {
		return false, nil
	}
	
	// Check if user is the creator
	if event.CreatedByID == userID {
		return true, nil
	}
	
	// Check if user is a participant
	participants, err := s.eventRepo.GetParticipants(eventID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve event participants: %w", err)
	}
	
	for _, participant := range participants {
		if participant.UserID == userID {
			return true, nil
		}
	}
	
	return false, nil
}

// CanModifyEvent checks if a user can modify an event (update, delete)
func (s *AuthorizationService) CanModifyEvent(ctx context.Context, userID, eventID uuid.UUID) (bool, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve event: %w", err)
	}
	
	// Only the creator can modify the event
	return event.CreatedByID == userID, nil
}

// CanCreateInviteCode checks if a user can create an invite code for an organization
func (s *AuthorizationService) CanCreateInviteCode(ctx context.Context, userID, orgID uuid.UUID) (bool, error) {
	// Only organization admins can create invite codes
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	// Check if the user is in the organization
	if user.OrganizationID == nil || *user.OrganizationID != orgID {
		return false, nil
	}
	
	// Check if the user is an admin
	return user.Role == RoleAdmin, nil
}

// CanManageOrganization checks if a user can manage an organization (update settings, etc.)
func (s *AuthorizationService) CanManageOrganization(ctx context.Context, userID, orgID uuid.UUID) (bool, error) {
	// Only organization admins can manage the organization
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	// Check if the user is in the organization
	if user.OrganizationID == nil || *user.OrganizationID != orgID {
		return false, nil
	}
	
	// Check if the user is an admin
	return user.Role == RoleAdmin, nil
}

// GetUserOrganizationID safely retrieves a user's organization ID
func (s *AuthorizationService) GetUserOrganizationID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	if user.OrganizationID == nil {
		return uuid.Nil, errors.New("user does not belong to an organization")
	}
	
	return *user.OrganizationID, nil
}

// ValidateUserAccess ensures a user belongs to a specific organization
func (s *AuthorizationService) ValidateUserAccess(ctx context.Context, userID, orgID uuid.UUID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	if user.OrganizationID == nil || *user.OrganizationID != orgID {
		return ErrNotMember
	}
	
	return nil
}