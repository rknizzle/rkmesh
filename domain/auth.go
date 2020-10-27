package domain

import (
	"context"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	SignUp(ctx context.Context, email string, password string) error
}
