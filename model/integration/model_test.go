package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
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

	var mockUserID int64 = 1
	// mock the request to have a JWT token containing a users ID
	c.Set("user", mockTokenWithUserID(mockUserID))

	err = mHandler.GetAll(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
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
