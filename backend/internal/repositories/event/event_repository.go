package event

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"hello-pulse.fr/internal/models/event"
	"hello-pulse.fr/internal/models/user"
)

// Repository handles database operations for events
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new event repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new event
func (r *Repository) Create(event *event.Event) error {
	return r.db.Create(event).Error
}

// FindByID finds an event by ID
func (r *Repository) FindByID(id uuid.UUID) (*event.Event, error) {
	var event event.Event
	err := r.db.First(&event, "event_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Update updates an event
func (r *Repository) Update(event *event.Event) error {
	return r.db.Save(event).Error
}

// Delete deletes an event
func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&event.Event{}, "event_id = ?", id).Error
}

// FindByOrganization returns events for a specific organization
func (r *Repository) FindByOrganization(orgID uuid.UUID) ([]event.Event, error) {
	var events []event.Event
	err := r.db.Where("organization_id = ?", orgID).Find(&events).Error
	return events, err
}

// FindByCreator returns events created by a specific user
func (r *Repository) FindByCreator(creatorID uuid.UUID) ([]event.Event, error) {
	var events []event.Event
	err := r.db.Where("created_by_id = ?", creatorID).Find(&events).Error
	return events, err
}

// AddParticipant adds a user to an event
func (r *Repository) AddParticipant(eventID, userID uuid.UUID) error {
	return r.db.Exec("INSERT INTO event_users (event_id, user_id) VALUES (?, ?)", eventID, userID).Error
}

// RemoveParticipant removes a user from an event
func (r *Repository) RemoveParticipant(eventID, userID uuid.UUID) error {
	return r.db.Exec("DELETE FROM event_users WHERE event_id = ? AND user_id = ?", eventID, userID).Error
}

// ClearParticipants removes all participants from an event
func (r *Repository) ClearParticipants(eventID uuid.UUID) error {
	return r.db.Exec("DELETE FROM event_users WHERE event_id = ?", eventID).Error
}

// GetParticipants gets all participants of an event
func (r *Repository) GetParticipants(eventID uuid.UUID) ([]user.User, error) {
	var users []user.User
	err := r.db.Joins("JOIN event_users ON users.user_id = event_users.user_id").
		Where("event_users.event_id = ?", eventID).
		Find(&users).Error
	return users, err
}

// FindForUser finds events for a specific user (created by or participating in)
func (r *Repository) FindForUser(userID uuid.UUID) ([]event.Event, error) {
	var events []event.Event
	err := r.db.
		Distinct("events.*").
		Joins("LEFT JOIN event_users ON events.event_id = event_users.event_id").
		Where("events.created_by_id = ? OR event_users.user_id = ?", userID, userID).
		Find(&events).Error
	return events, err
}

// UpdateTitle updates the title of an event
func (r *Repository) UpdateTitle(eventID uuid.UUID, title string) error {
	return r.db.Model(&event.Event{}).Where("event_id = ?", eventID).Update("title", title).Error
}

// FindUpcoming finds upcoming events after a specific date
func (r *Repository) FindUpcoming(orgID uuid.UUID, after time.Time) ([]event.Event, error) {
	var events []event.Event
	err := r.db.Where("organization_id = ? AND date >= ?", orgID, after).
		Order("date ASC, start_time ASC").
		Find(&events).Error
	return events, err
}