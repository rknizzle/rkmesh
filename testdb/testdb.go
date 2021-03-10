package testdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rknizzle/rkmesh/domain"
)

const (
	dbUser = "postgres"
	dbPass = "postgres"
	dbName = "postgres"
	dbHost = "localhost"
	dbPort = "5432"
)

type TestDB struct {
	Conn *sql.DB
}

// Open creates a connection to a test database and returns an instance of TestDB which can be used
// to truncate and seed data into it. It also returns a connection to the test database that can be
// used by the application while running integration tests.
func Open() (TestDB, *sql.DB, error) {
	connection := fmt.Sprintf(
		`host=%s port=%s user=%s
		password=%s dbname=%s sslmode=disable`,
		dbHost, dbPort, dbUser, dbPass, dbName)

	dbConn, err := sql.Open("postgres", connection)
	if err != nil {
		return TestDB{}, nil, err
	}

	tdb := TestDB{Conn: dbConn}
	return tdb, dbConn, nil
}

// Truncate removes all seed data from the test database
func (t *TestDB) Truncate() error {
	query := "TRUNCATE TABLE models, users;"

	stmt, err := t.Conn.PrepareContext(context.TODO(), query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

// SeedModels places test data into the test database for integration tests
func (t *TestDB) SeedModels() ([]domain.Model, error) {
	testModels := testModels()

	for _, m := range testModels {
		query := `INSERT INTO models (name, user_id, updated_at, created_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
		stmt, err := t.Conn.Prepare(query)
		if err != nil {
			return nil, err
		}

		var ID int64
		err = stmt.QueryRow(m.Name, m.UserID).Scan(&ID)
		if err != nil {
			return nil, err
		}

		m.ID = ID
	}

	return testModels, nil
}

// SeedUsers places test users into the test database for integration tests
func (t *TestDB) SeedUsers() ([]domain.User, error) {
	testUsers := testUsers()

	for _, u := range testUsers {
		query := `INSERT INTO users (id, email, password, updated_at, created_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
		stmt, err := t.Conn.Prepare(query)
		if err != nil {
			return nil, err
		}

		var ID int64
		err = stmt.QueryRow(u.ID, u.Email, u.Password).Scan(&ID)
		if err != nil {
			return nil, err
		}

		u.ID = ID
	}

	return testUsers, nil
}

func testModels() []domain.Model {
	models := []domain.Model{
		{
			Name:   "test.stl",
			UserID: 1,
		},
		{
			Name:   "test2.stl",
			UserID: 1,
		},
		{
			Name:   "test3.stl",
			UserID: 2,
		},
	}
	return models
}

func testUsers() []domain.User {
	users := []domain.User{
		{
			ID:       1,
			Email:    "example@gmail.com",
			Password: "password",
		},
		{
			ID:       2,
			Email:    "example2@gmail.com",
			Password: "password",
		},
	}
	return users
}
