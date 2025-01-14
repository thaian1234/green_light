package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thaian1234/green_light/internal/core/domain"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (r *MovieRepository) Insert(ctx context.Context, movie *domain.Movie) error {
	query := `
        INSERT INTO movies (title, year, runtime, genres)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version
    `
	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		movie.Genres,
	}

	return r.db.QueryRow(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (r *MovieRepository) GetByID(ctx context.Context, id int64) (*domain.Movie, error) {
	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1
	`
	var movie domain.Movie
	err := r.db.QueryRow(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		&movie.Genres,
		&movie.Version,
	)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) Update(ctx context.Context, movie *domain.Movie) error {
	return nil
}

func (r *MovieRepository) Delete(ctx context.Context, id int64) error {
	return nil
}
