package service

import (
	"context"
	"time"

	"github.com/rknizzle/rkmesh/domain"
)

type authService struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewAuthService(u domain.UserRepository, timeout time.Duration) domain.AuthService {
	return &authService{
		userRepo:       u,
		contextTimeout: timeout,
	}
}

func (a *authService) Login(c context.Context, email string, password string) (token string, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	_, err = a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// do a hash verify on user.password and check if it matches the one in the database
	//hash.verify(password, user.password)

	// generate a token if it is valid and return the token
	return "xxx", nil
}

func (a *authService) SignUp(c context.Context, email string, password string) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	// hash the password before placing the user into the database
	// hashPassword()

	err = a.userRepo.Create(ctx, &domain.User{Email: email, Password: password})
	if err != nil {
		return err
	}

	return
}
