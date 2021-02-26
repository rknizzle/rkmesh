package model_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/rknizzle/rkmesh/domain"

	"github.com/rknizzle/rkmesh/model"
)

func TestPGGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockModels := []domain.Model{
		domain.Model{
			ID: 1, Name: "test1.stl", UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		domain.Model{
			ID: 2, Name: "test2.stl", UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "updated_at", "created_at"}).
		AddRow(mockModels[0].ID, mockModels[0].Name, mockModels[0].UpdatedAt, mockModels[0].CreatedAt).
		AddRow(mockModels[1].ID, mockModels[1].Name, mockModels[1].UpdatedAt, mockModels[1].CreatedAt)

	query := `[SELECT * FROM models ORDER BY created_at]`

	mock.ExpectQuery(query).WillReturnRows(rows)

	m := model.NewPostgresModelRepository(db)

	list, err := m.GetAll(context.TODO())

	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestPGGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "updated_at", "created_at"}).
		AddRow(1, "test1.stl", time.Now(), time.Now())

	query := "[SELECT * FROM models WHERE id = $1]"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := model.NewPostgresModelRepository(db)

	num := int64(5)
	model, err := a.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, model)
}

func TestPGStore(t *testing.T) {
	m := &domain.Model{
		Name: "test.stl", UpdatedAt: time.Now(), CreatedAt: time.Now(),
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "[INSERT INTO models (name, updated_at, created_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id]"
	prep := mock.ExpectPrepare(query)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	prep.ExpectQuery().WithArgs(m.Name).WillReturnRows(rows)

	p := model.NewPostgresModelRepository(db)

	err = p.Store(context.TODO(), m)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), m.ID)
}

func TestPGGetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "updated_at", "created_at"}).
		AddRow(1, "test1.stl", time.Now(), time.Now())

	query := `[SELECT * FROM models WHERE name = '$1']`

	mock.ExpectQuery(query).WillReturnRows(rows)
	p := model.NewPostgresModelRepository(db)

	name := "test.stl"
	model, err := p.GetByName(context.TODO(), name)
	assert.NoError(t, err)
	assert.NotNil(t, model)
}

func TestPGDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "[DELETE FROM models WHERE id = $1]"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	p := model.NewPostgresModelRepository(db)

	num := int64(12)
	err = p.Delete(context.TODO(), num)
	assert.NoError(t, err)
}
