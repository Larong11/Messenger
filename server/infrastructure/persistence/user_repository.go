package persistence

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"server/domain/user"
	"time"
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
func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *user.User) (int, error) {
	query := `INSERT INTO users (
                    first_name, last_name, username, email, password_hash, is_email_verified, created_at, avatar_url, 
                   last_seen_at, user_status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;`
	row := r.pool.QueryRow(ctx, query,
		user.FirstName, user.LastName, user.UserName, user.Email, user.PasswordHash, user.IsEmailVerified, time.Now().UTC(),
		"url", time.Now().UTC(), user.UserStatus)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
