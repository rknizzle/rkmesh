package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rknizzle/rkmesh/auth"
	"github.com/rknizzle/rkmesh/domain/mocks"
)

func TestLogin(t *testing.T) {
	mockService := new(mocks.AuthService)

	mockService.On("Login", mock.Anything, "ryan@example.com", "password").Return("token goes here", nil)

	tempLoginInput := &auth.UserInput{Email: "ryan@example.com", Password: "password"}
	j, err := json.Marshal(tempLoginInput)
	assert.NoError(t, err)

	e := echo.New()

	req, err := http.NewRequest(echo.POST, "/auth/login", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("auth/login")

	handler := auth.AuthHandler{
		Service: mockService,
	}
	err = handler.Login(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestSignUp(t *testing.T) {
	mockService := new(mocks.AuthService)

	mockService.On("SignUp", mock.Anything, "ryan@example.com", "password").Return(nil)

	tempLoginInput := &auth.UserInput{Email: "ryan@example.com", Password: "password"}
	j, err := json.Marshal(tempLoginInput)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/auth/sign-up", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/auth/sign-up")

	handler := auth.AuthHandler{
		Service: mockService,
	}
	err = handler.SignUp(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}
