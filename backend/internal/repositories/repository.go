package repositories

import "github.com/google/uuid"

// Repository is a generic repository interface
type Repository[T any] interface {
	Create(model *T) error
	FindByID(id uuid.UUID) (*T, error)
	Update(model *T) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]T, error)
}