package model

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/rknizzle/rkmesh/domain"
)

type responseError struct {
	Message string `json:"message"`
}

type ModelHandler struct {
	Service domain.ModelService
}

// NewModelHandler will initialize the /models resources endpoints
func NewModelHandler(e *echo.Group, s domain.ModelService) {
	handler := &ModelHandler{
		Service: s,
	}

	// /models...
	e.GET("", handler.GetAll)
	e.POST("", handler.Store)
	e.GET("/:id", handler.GetByID)
	e.DELETE("/:id", handler.Delete)
}

func (m *ModelHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	userID := getUserIDFromRequest(c)

	mList, err := m.Service.GetAllUserModels(ctx, userID)
	if err != nil {
		return c.JSON(getStatusCode(err), responseError{Message: err.Error()})
	}

	//c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, mList)
}

func (m *ModelHandler) GetByID(c echo.Context) error {
	// convert the url param 'id' from a string to int64
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	ctx := c.Request().Context()
	userID := getUserIDFromRequest(c)

	model, err := m.Service.GetByID(ctx, id, userID)
	if err != nil {
		return c.JSON(getStatusCode(err), responseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model)
}

func isRequestValid(m *domain.Model) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *ModelHandler) Store(c echo.Context) (err error) {
	userID := getUserIDFromRequest(c)

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseError{Message: err.Error()})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(getStatusCode(err), responseError{Message: err.Error()})
	}
	defer src.Close()

	model := &domain.Model{}

	ctx := c.Request().Context()
	err = m.Service.Store(ctx, model, src, file.Filename, userID)
	if err != nil {
		return c.JSON(getStatusCode(err), responseError{Message: err.Error()})
	}

	// TODO: only return the values that I want returned to the user

	return c.JSON(http.StatusCreated, model)
}

// Delete will delete model by given param
func (m *ModelHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()
	userID := getUserIDFromRequest(c)

	err = m.Service.Delete(ctx, id, userID)
	if err != nil {
		return c.JSON(getStatusCode(err), responseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
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

func getUserIDFromRequest(c echo.Context) int64 {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return int64(claims["user_id"].(float64))
}
