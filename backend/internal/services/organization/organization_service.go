package organization

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/invite"
	"hello-pulse.fr/internal/models/organization"
	"hello-pulse.fr/internal/models/user"
	inviterepo "hello-pulse.fr/internal/repositories/invite"
	orgrepo "hello-pulse.fr/internal/repositories/organization"
	userrepo "hello-pulse.fr/internal/repositories/user"
)

// Service handles organization business logic
type Service struct {
	orgRepo    *orgrepo.Repository
	userRepo   *userrepo.Repository
	inviteRepo *inviterepo.Repository
}

// NewService creates a new organization service
func NewService(orgRepo *orgrepo.Repository, userRepo *userrepo.Repository, inviteRepo *inviterepo.Repository) *Service {
	return &Service{
		orgRepo:    orgRepo,
		userRepo:   userRepo,
		inviteRepo: inviteRepo,
	}
}

// GetInviteCodes gets all invite codes for an organization
func (s *Service) GetInviteCodes(orgID uuid.UUID) ([]invite.InviteCode, error) {
    return s.inviteRepo.FindByOrganization(orgID)
}

// CreateOrganization creates a new organization
func (s *Service) CreateOrganization(name string, ownerID uuid.UUID) (*organization.Organization, error) {
	// Check if organization name already exists
	existingOrg, err := s.orgRepo.FindByName(name)
	if err == nil && existingOrg != nil {
		return nil, errors.New("organization name already exists")
	}

	// Get the owner user
	owner, err := s.userRepo.FindByID(ownerID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Create organization
	org := &organization.Organization{
		OrganizationName: name,
		OwnerID:          ownerID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.orgRepo.Create(org); err != nil {
		return nil, err
	}

	// Update user's organization and role
	owner.OrganizationID = &org.OrganizationID
	owner.Role = "Admin"
	if err := s.userRepo.Update(owner); err != nil {
		return nil, err
	}

	return org, nil
}

// CreateInviteCode creates a new invite code for an organization
func (s *Service) CreateInviteCode(orgID uuid.UUID, expirationMs int64) (*invite.InviteCode, error) {
	// Check if organization exists
	org, err := s.orgRepo.FindByID(orgID)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	// Generate invite code
	code, err := s.generateRandomString(6)
	if err != nil {
		return nil, err
	}

	// Create invite code
	inviteCode := &invite.InviteCode{
		Value:          code,
		OrganizationID: org.OrganizationID,
		ExpirationTime: s.msToTime(expirationMs),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.inviteRepo.Create(inviteCode); err != nil {
		return nil, err
	}

	return inviteCode, nil
}

// JoinOrganization adds a user to an organization using an invite code
func (s *Service) JoinOrganization(userID uuid.UUID, code string) error {
	// Get invite code
	inviteCode, err := s.inviteRepo.FindByCode(code)
	if err != nil {
		return errors.New("invalid invite code")
	}

	// Check if invite code is expired
	if inviteCode.ExpirationTime.Before(time.Now()) {
		return errors.New("invite code expired")
	}

	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Update user's organization and role
	user.OrganizationID = &inviteCode.OrganizationID
	user.Role = "User"
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}

// GetOrganizationUsers gets all users in an organization
func (s *Service) GetOrganizationUsers(orgID uuid.UUID) ([]user.User, error) {
	return s.userRepo.FindByOrganization(orgID)
}

// generateRandomString generates a random string of the specified length
func (s *Service) generateRandomString(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[randomInt.Int64()]
	}

	return string(result), nil
}

// msToTime converts milliseconds to time.Time
func (s *Service) msToTime(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond))
}