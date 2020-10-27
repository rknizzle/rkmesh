package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"github.com/rknizzle/rkmesh/domain"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// AuthHandler  represent the httphandler for auth
type AuthHandler struct {
	Service domain.AuthService
}

// NewAuthHandler will initialize the auth/ resources endpoint
func NewAuthHandler(e *echo.Echo, s domain.AuthService) {
	handler := &AuthHandler{
		Service: s,
	}

	e.POST("/auth/sign-up", handler.SignUp)
	e.POST("/auth/login", handler.Login)
}

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (a *AuthHandler) SignUp(c echo.Context) error {
	var input UserInput
	err := c.Bind(&input)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	err = a.Service.SignUp(ctx, input.Email, input.Password)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (a *AuthHandler) Login(c echo.Context) error {
	var input UserInput
	err := c.Bind(&input)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	token, err := a.Service.Login(ctx, input.Email, input.Password)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, LoginResponse{Token: token})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
