package persistence

import (
	"context"
	"errors"
	"server/domain/user"
	upgradeerrors "server/internal/errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		return nil, upgradeerrors.NewInternal("db error") // другая ошибка
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
		return nil, upgradeerrors.NewInternal("db error")
	}
	return &u.ID, nil
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *user.User) (int, error) {
	query := `INSERT INTO users (
	        first_name, last_name, username, email, password_hash, is_email_verified,
	        created_at, avatar_url, last_seen_at, user_status
	    ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id;`

	var id int
	err := r.pool.QueryRow(ctx, query,
		user.FirstName, user.LastName, user.UserName, user.Email, user.PasswordHash,
		user.IsEmailVerified, time.Now().UTC(), "url", time.Now().UTC(), user.UserStatus,
	).Scan(&id)
	if err != nil {
		return -1, upgradeerrors.NewInternal("db error")
	}
	return id, nil
}

func (r *PostgresUserRepository) CreateVerificationCode(ctx context.Context, email string, verificationCode string) error {
	queryCode := `INSERT INTO user_verification_codes (email, verification_code, created_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (email) DO UPDATE SET
	  verification_code = EXCLUDED.verification_code,
	  created_at = EXCLUDED.created_at;`
	_, err := r.pool.Exec(ctx, queryCode, email, verificationCode, time.Now().UTC())
	if err != nil {
		return upgradeerrors.NewInternal("db error")
	}
	return nil
}
func (r *PostgresUserRepository) GetVerificationCode(ctx context.Context, email string) (string, time.Time, error) {
	queryCode := `SELECT verification_code, created_at FROM user_verification_codes WHERE email = $1;`
	row := r.pool.QueryRow(ctx, queryCode, email)
	var verificationCode string
	var createdAt time.Time
	err := row.Scan(&verificationCode, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", time.Time{}, upgradeerrors.NewBadRequest("verification code not found")
		}
		return "", time.Time{}, upgradeerrors.NewInternal("db error")
	}
	return verificationCode, createdAt, nil
}
