package persistence

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"server/domain/user"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}
func (r *PostgresUserRepository) FindByUserName(ctx context.Context, userName string) (*int, error) {
	query := `SELECT id FROM users WHERE username = $1;`

	row := r.pool.QueryRow(ctx, query, userName)

	var u user.User
	err := row.Scan(&u.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // пользователь не найден
		}
		return nil, err // другая ошибка
	}

	return &u.ID, nil

}
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*int, error) {
	query := `SELECT id FROM users WHERE email = $1;`
	row := r.pool.QueryRow(ctx, query, email)
	var u user.User
	err := row.Scan(&u.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u.ID, nil
}
