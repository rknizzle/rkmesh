package domain

import (
	"context"
	"time"
)

type Model struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Volume      float64   `json:"volume"`
	SurfaceArea float64   `json:"surface_area"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// ModelService represent the models business logic
type ModelService interface {
	GetAll(ctx context.Context) ([]Model, error)
	GetByID(ctx context.Context, id int64) (Model, error)
	GetByName(ctx context.Context, name string) (Model, error)
	Store(context.Context, *Model) error
	Delete(ctx context.Context, id int64) error
}

// ModelService represent the models repository contract
type ModelRepository interface {
	GetAll(ctx context.Context) (res []Model, err error)
	GetByID(ctx context.Context, id int64) (Model, error)
	GetByName(ctx context.Context, name string) (Model, error)
	Store(ctx context.Context, m *Model) error
	Delete(ctx context.Context, id int64) error
}
