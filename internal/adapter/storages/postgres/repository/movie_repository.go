package repository

import (
	"context"
	"fmt"
	"strings"

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

	err := r.db.QueryRow(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return domain.ErrConflictingData
			case "foreign_key_violation":
				return domain.ErrDataNotFound
			}
		}
		return domain.ErrInternalServer
	}
	return nil
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

func (r *MovieRepository) GetAll(ctx context.Context, title string, genres []string, filter domain.Filter) ([]*domain.Movie, domain.Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filter.SortColumn(), filter.SortDirection())
	args := []any{
		strings.TrimSpace(strings.ToLower(title)),
		pq.Array(genres),
		filter.Limit(),
		filter.Offset(),
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, domain.Metadata{}, domain.ErrInternalServer
	}
	defer rows.Close()

	totalRecords := 0
	movies := make([]*domain.Movie, 0)
	for rows.Next() {
		var movie domain.Movie
		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			&movie.Genres,
			&movie.Version,
		)
		if err != nil {
			return nil, domain.Metadata{}, domain.ErrInternalServer
		}
		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.Metadata{}, domain.ErrInternalServer
	}
	metadata := domain.CalculateMetadata(totalRecords, filter.Page, filter.Size)
	return movies, metadata, nil
}

func (r *MovieRepository) Update(ctx context.Context, movie *domain.Movie) error {
	query := `
        UPDATE movies
        SET title = COALESCE($1, title),
            year = COALESCE($2, year),
            runtime = COALESCE($3, runtime),
            genres = COALESCE($4, genres),
            version = version + 1
        WHERE id = $5 and version = $6
        RETURNING version
    `
	args := []any{
		util.NullString(movie.Title),
		util.NullUint64(uint64(movie.Year)),
		util.NullUint64(uint64(movie.Runtime)),
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}
	err := r.db.QueryRow(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrUpdateConflict
		}
		return domain.ErrInternalServer
	}
	return nil
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
