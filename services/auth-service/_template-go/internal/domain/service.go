package domain

import (
	"context"
	"time"
)

// Service interface - define your business logic here
type ServiceInterface interface {
	// Add your service methods here
}

// Example entity - replace with your domain model
type Entity struct {
	ID        string    `db:"id,primarykey"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Example service - replace with your business logic
type Service struct {
	repo Repository
}

// Repository interface - implement in infrastructure layer
type Repository interface {
	Create(ctx context.Context, entity *Entity) error
	GetByID(ctx context.Context, id string) (*Entity, error)
	Update(ctx context.Context, entity *Entity) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Entity, error)
}

// NewService creates a new service instance
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Example method - implement your business logic
func (s *Service) GetEntity(ctx context.Context, id string) (*Entity, error) {
	return s.repo.GetByID(ctx, id)
}
