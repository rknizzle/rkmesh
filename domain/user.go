package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository represent the users repository contract
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, u *User) error
}
