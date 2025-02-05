package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thaian1234/green_light/internal/core/domain"
	"github.com/thaian1234/green_light/pkg/logger"
	"github.com/thaian1234/green_light/pkg/util"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Insert(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
	`
	args := []any{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
	}
	logger.Debug("Inserting user", "args", args)
	err := r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		logger.Debug("Inserting user", "err", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505" && strings.Contains(pgErr.Message, "users_email_key"):
				return domain.ErrDuplicatedEmail
			}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrDataNotFound
		}
		return domain.ErrInternalServer
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, created_at, name, email, password_hash, activated, version
		FROM users
		WHERE email = $1
	`
	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case err == pgx.ErrNoRows:
			return nil, domain.ErrDataNotFound
		default:
			return nil, domain.ErrInternalServer
		}
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = COALESCE($1, name),
			email = COALESCE($2, email),
			password_hash = COALESCE($3, password_hash),
			activated = COALESCE($4, activated),
			version = version + 1
		WHERE id = $5 and version = $6
		RETURNING version
	`
	args := []any{
		util.NullString(user.Name),
		util.NullString(user.Email),
		util.NullString(string(user.Password.Hash)),
		user.Activated,
		user.ID,
		user.Version,
	}
	err := r.db.QueryRow(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err == pgx.ErrNoRows:
			return domain.ErrUpdateConflict
		default:
			return domain.ErrInternalServer
		}
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM users WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return domain.ErrInternalServer
	}
	if result.RowsAffected() == 0 {
		return domain.ErrDataNotFound
	}
	return nil
}
