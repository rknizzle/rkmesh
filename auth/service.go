package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rknizzle/rkmesh/domain"
	"golang.org/x/crypto/bcrypt"
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

	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("Unauthorized")
	}

	token, err = CreateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *authService) SignUp(c context.Context, email string, password string) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	// hash the password before placing the user into the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return err
	}

	err = a.userRepo.Create(ctx, &domain.User{Email: email, Password: string(hashedPassword)})
	if err != nil {
		return err
	}

	return
}

func CreateToken(userid int64) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return token, nil
}
