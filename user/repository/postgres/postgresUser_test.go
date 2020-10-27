package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/rknizzle/rkmesh/domain"
	userPostgresRepo "github.com/rknizzle/rkmesh/user/repository/postgres"
)

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "updated_at", "created_at"}).
		AddRow(1, "ryan@example.com", "password", time.Now(), time.Now())

	query := "[SELECT * FROM users WHERE id = $1]"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := userPostgresRepo.NewPostgresUserRepository(db)

	user, err := a.GetByEmail(context.TODO(), "ryan@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestCreate(t *testing.T) {
	u := &domain.User{
		Email: "ryan@example.com", Password: "password", UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "[INSERT INTO users (email, password, updated_at, created_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id]"
	prep := mock.ExpectPrepare(query)

	// this is what its mocked to create after running userRepo.Create() (AKA its
	// just gonna return a row with an id of 1)
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	prep.ExpectQuery().WithArgs(u.Email, u.Password).WillReturnRows(rows)

	p := userPostgresRepo.NewPostgresUserRepository(db)

	err = p.Create(context.TODO(), u)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), u.ID)
}
