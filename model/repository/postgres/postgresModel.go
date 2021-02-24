package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/rknizzle/rkmesh/domain"
)

type postgresModelRepository struct {
	Conn *sql.DB
}

// NewPostgresModelRepository will create an object that represent the model.Repository interface
func NewPostgresModelRepository(Conn *sql.DB) domain.ModelRepository {
	return &postgresModelRepository{Conn}
}

// gets all rows from the result of a sql query
func (p *postgresModelRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Model, err error) {
	rows, err := p.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Model, 0)
	for rows.Next() {
		t := domain.Model{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (p *postgresModelRepository) GetAll(ctx context.Context) (res []domain.Model, err error) {
	query := `SELECT * FROM models ORDER BY created_at`

	res, err = p.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return
}

func (p *postgresModelRepository) GetByID(ctx context.Context, id int64) (res domain.Model, err error) {
	query := `SELECT * FROM models WHERE id = $1`

	list, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Model{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (p *postgresModelRepository) GetByName(ctx context.Context, name string) (res domain.Model, err error) {
	query := `SELECT * FROM models WHERE name = '$1'`

	list, err := p.fetch(ctx, query, name)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (p *postgresModelRepository) Store(ctx context.Context, m *domain.Model) (err error) {
	query := `INSERT INTO models (name, updated_at, created_at) VALUES ($1, NOW(), NOW()) RETURNING id`
	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	var ID int64
	err = stmt.QueryRowContext(ctx, m.Name).Scan(&ID)
	if err != nil {
		return
	}

	m.ID = ID
	return
}

func (p *postgresModelRepository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM models WHERE id = $1`

	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}
