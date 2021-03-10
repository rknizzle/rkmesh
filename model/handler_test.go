package model_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rknizzle/rkmesh/domain"
	"github.com/rknizzle/rkmesh/domain/mocks"
	"github.com/rknizzle/rkmesh/model"
)

func TestHandlerGetAll(t *testing.T) {
	mockModelList := []domain.Model{
		{
			ID:     55,
			Name:   "test.stl",
			UserID: 1,
		},
		{
			ID:     56,
			Name:   "test2.stl",
			UserID: 1,
		},
	}

	mockService := new(mocks.ModelService)

	var mockUserID int64 = 1
	mockService.On("GetAllUserModels", mock.Anything, mockUserID).Return(mockModelList, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/models", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", mockTokenWithUserID(mockUserID))

	handler := model.ModelHandler{
		Service: mockService,
	}
	err = handler.GetAll(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandlerGetAllError(t *testing.T) {
	mockService := new(mocks.ModelService)

	var mockUserID int64 = 1
	mockService.On("GetAllUserModels", mock.Anything, mockUserID).Return(nil, domain.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/models", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", mockTokenWithUserID(mockUserID))

	handler := model.ModelHandler{
		Service: mockService,
	}
	err = handler.GetAll(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandlerGetByID(t *testing.T) {
	var mockModel domain.Model
	err := faker.FakeData(&mockModel)
	assert.NoError(t, err)

	mockService := new(mocks.ModelService)

	num := int(mockModel.ID)
	var mockUserID int64 = 1

	mockService.On("GetByID", mock.Anything, int64(num), mockUserID).Return(mockModel, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/models/"+strconv.Itoa(num), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("model/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	c.Set("user", mockTokenWithUserID(mockUserID))

	handler := model.ModelHandler{
		Service: mockService,
	}

	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandlerStore(t *testing.T) {
	var mockModel domain.Model
	err := faker.FakeData(&mockModel)
	assert.NoError(t, err)

	mockService := new(mocks.ModelService)

	var mockUserID int64 = 1

	mockService.On("Store", mock.Anything, mock.AnythingOfType("*domain.Model"), mock.Anything, "test.stl", mockUserID).Return(nil)

	e := echo.New()
	formData, multipartBoundary, err := mockFormData()
	assert.NoError(t, err)

	req, err := http.NewRequest(echo.POST, "/models", formData)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", multipartBoundary)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// mock the request to have a JWT token containing a users ID
	c.Set("user", mockTokenWithUserID(mockUserID))
	c.SetPath("/models")

	handler := model.ModelHandler{
		Service: mockService,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandlerDelete(t *testing.T) {
	var mockModel domain.Model
	err := faker.FakeData(&mockModel)
	assert.NoError(t, err)

	mockService := new(mocks.ModelService)
	var mockUserID int64 = 1

	num := int(mockModel.ID)

	mockService.On("Delete", mock.Anything, int64(num), mockUserID).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/models/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("models/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	c.Set("user", mockTokenWithUserID(mockUserID))

	handler := model.ModelHandler{
		Service: mockService,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}

func mockFormData() (*bytes.Buffer, string, error) {
	b := new(bytes.Buffer)
	writer := multipart.NewWriter(b)
	part, err := writer.CreateFormFile("file", "test.stl")
	if err != nil {
		return &bytes.Buffer{}, "", err
	}
	part.Write([]byte("file data here"))

	err = writer.Close()
	if err != nil {
		return &bytes.Buffer{}, "", err
	}

	return b, writer.FormDataContentType(), nil
}

func mockTokenWithUserID(mockUserID int64) *jwt.Token {

	// NOTE: for some reason Echo's JWT middleware has the user_id as a float64 value so im mocking
	// here to have a float64 user_id value to match the input that Echo will give me
	mockUserIDAsFloat64 := float64(mockUserID)

	// mock the JWT token to be from a user with an ID of '1'
	mockToken := &jwt.Token{
		Claims: jwt.MapClaims{
			"user_id": mockUserIDAsFloat64,
		},
	}

	return mockToken
}
