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
	query := "TRUNCATE TABLE models;"

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
	models := []domain.Model{
		{
			Name: "test.stl",
		},
		{
			Name: "test2.stl",
		},
		{
			Name: "test3.stl",
		},
	}

	for _, m := range models {
		query := `INSERT INTO models (name, updated_at, created_at) VALUES ($1, NOW(), NOW()) RETURNING id`
		stmt, err := t.Conn.Prepare(query)
		if err != nil {
			return nil, err
		}

		var ID int64
		err = stmt.QueryRow(m.Name).Scan(&ID)
		if err != nil {
			return nil, err
		}

		m.ID = ID
	}

	return models, nil
}
