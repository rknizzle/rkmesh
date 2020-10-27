package postgres

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/rknizzle/rkmesh/domain"
)

type postgresUserRepository struct {
	Conn *sql.DB
}

// NewPostgresUserRepository will create an object that represent the user.Repository interface
func NewPostgresUserRepository(Conn *sql.DB) domain.UserRepository {
	return &postgresUserRepository{Conn}
}

func (u *postgresUserRepository) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
	query := `SELECT * FROM users WHERE email = $1`

	list, err := u.fetch(ctx, query, email)
	if err != nil {
		return domain.User{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

// gets all rows from the result of a sql query
func (p *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {

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

	result = make([]domain.User, 0)
	for rows.Next() {
		t := domain.User{}
		err = rows.Scan(
			&t.ID,
			&t.Email,
			&t.Password,
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

func (p *postgresUserRepository) Create(ctx context.Context, u *domain.User) (err error) {
	query := `INSERT INTO users (email, password, updated_at, created_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	var ID int64
	err = stmt.QueryRowContext(ctx, u.Email, u.Password).Scan(&ID)
	if err != nil {
		return
	}

	u.ID = ID
	return
}
