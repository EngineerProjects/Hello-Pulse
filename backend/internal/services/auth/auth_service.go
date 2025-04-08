package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	authrepo "hello-pulse.fr/internal/repositories/auth" 
	"hello-pulse.fr/internal/models/auth"
	"hello-pulse.fr/internal/models/user"
	userrepo "hello-pulse.fr/internal/repositories/user"
)

// Service handles authentication business logic
type Service struct {
	userRepo    *userrepo.Repository
	sessionRepo *authrepo.Repository 
}

// NewService creates a new authentication service
func NewService(userRepo *userrepo.Repository, sessionRepo *authrepo.Repository) *Service {
	return &Service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// Rest of the code remains the same
// RegisterUser registers a new user
func (s *Service) RegisterUser(firstName, lastName, email, password, phone, address string) (*user.User, string, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, "", errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	newUser := &user.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Phone:        phone,
		Address:      address,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, "", err
	}

	// Create session
	token := uuid.New().String()
	session := &auth.Session{
		UserID:    newUser.UserID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, "", err
	}

	return newUser, token, nil
}

// Login authenticates a user and returns a session token
func (s *Service) Login(email, password string) (*user.User, string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Create session
	token := uuid.New().String()
	session := &auth.Session{
		UserID:    user.UserID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, "", err
	}

	// Update last active
	user.LastActive = time.Now()
	if err := s.userRepo.Update(user); err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Logout invalidates a user session
func (s *Service) Logout(token string) error {
	return s.sessionRepo.DeleteByToken(token)
}

// ValidateSession checks if a session is valid
func (s *Service) ValidateSession(token string) (*user.User, error) {
	// Find session
	session, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return nil, errors.New("invalid session")
	}

	// Check if session is expired
	if session.ExpiresAt.Before(time.Now()) {
		_ = s.sessionRepo.Delete(session.SessionID)
		return nil, errors.New("session expired")
	}

	// Get user
	user, err := s.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}