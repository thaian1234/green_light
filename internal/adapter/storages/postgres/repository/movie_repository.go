package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/thaian1234/green_light/internal/core/domain"
	"github.com/thaian1234/green_light/pkg/util"
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
		pq.Array(movie.Genres),
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
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.ErrInternalServer
	}
	return &movie, nil
}

func (r *MovieRepository) GetAll(ctx context.Context) ([]*domain.Movie, error) {
	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		ORDER BY id
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	defer rows.Close()
	movies := make([]*domain.Movie, 0)
	for rows.Next() {
		var movie domain.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			&movie.Genres,
			&movie.Version,
		)
		if err != nil {
			return nil, domain.ErrInternalServer
		}
		movies = append(movies, &movie)
	}
	if err := rows.Err(); err != nil {
		return nil, domain.ErrInternalServer
	}
	return movies, nil
}

func (r *MovieRepository) Update(ctx context.Context, movie *domain.Movie) error {
	query := `
        UPDATE movies
        SET title = COALESCE($1, title),
            year = COALESCE($2, year),
            runtime = COALESCE($3, runtime),
            genres = COALESCE($4, genres),
            version = version + 1
        WHERE id = $5 
        RETURNING version
    `
	args := []any{
		util.NullString(movie.Title),
		util.NullUint64(uint64(movie.Year)),
		util.NullUint64(uint64(movie.Runtime)),
		pq.Array(movie.Genres),
		movie.ID,
	}
	return r.db.QueryRow(ctx, query, args...).Scan(&movie.Version)
}

func (r *MovieRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM movies WHERE id = $1
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
