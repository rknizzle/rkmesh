package http_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rknizzle/rkmesh/domain"
	"github.com/rknizzle/rkmesh/domain/mocks"
	modelHTTP "github.com/rknizzle/rkmesh/model/controller/http"
)

func TestGetAll(t *testing.T) {
	var mockModel domain.Model
	err := faker.FakeData(&mockModel)
	assert.NoError(t, err)
	mockService := new(mocks.ModelService)
	mockListModel := make([]domain.Model, 0)
	mockListModel = append(mockListModel, mockModel)
	mockService.On("GetAll", mock.Anything).Return(mockListModel, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/models", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := modelHTTP.ModelHandler{
		Service: mockService,
	}
	err = handler.GetAll(c)
	require.NoError(t, err)

	// TODO: can I verify that it returns an array of Models?
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllError(t *testing.T) {
	mockService := new(mocks.ModelService)
	mockService.On("GetAll", mock.Anything).Return(nil, domain.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/models", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := modelHTTP.ModelHandler{
		Service: mockService,
	}
	err = handler.GetAll(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockModel domain.Model
	err := faker.FakeData(&mockModel)
	assert.NoError(t, err)

	mockService := new(mocks.ModelService)

	num := int(mockModel.ID)

	mockService.On("GetByID", mock.Anything, int64(num)).Return(mockModel, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/models/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("model/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := modelHTTP.ModelHandler{
		Service: mockService,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockModel := domain.Model{
		Name:      "test.stl",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockModel := mockModel
	tempMockModel.ID = 0
	mockService := new(mocks.ModelService)

	mockService.On("Store", mock.Anything, mock.AnythingOfType("*domain.Model"), mock.Anything, "test.stl").Return(nil)

	e := echo.New()
	formData, multipartBoundary, err := mockFormData()
	assert.NoError(t, err)

	req, err := http.NewRequest(echo.POST, "/models", formData)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", multipartBoundary)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/models")

	handler := modelHTTP.ModelHandler{
		Service: mockService,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockModel domain.Model
	err := faker.FakeData(&mockModel)
	assert.NoError(t, err)

	mockService := new(mocks.ModelService)

	num := int(mockModel.ID)

	mockService.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/models/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("models/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := modelHTTP.ModelHandler{
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
