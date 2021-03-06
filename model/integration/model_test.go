package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/rknizzle/rkmesh/model"
	"github.com/rknizzle/rkmesh/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Global variables that are defined in main_test.go and necessary for running the tests

// Model handler for passing the test requests into the application
var mHandler model.ModelHandler

// Test database for seeding and clearing data to fabricate particular situations for each test
var tdb testdb.TestDB

// TODO: not in use yet
//var tfs testFileStore.TestFileStore

// GET /models
func TestGetAll(t *testing.T) {
	tdb.Truncate()
	_, err := tdb.SeedModels()
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/models", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = mHandler.GetAll(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
}
