package ports

import (
	"context"

	"github.com/thaian1234/green_light/internal/core/domain"
)

type MovieRepository interface {
	Insert(ctx context.Context, movie *domain.Movie) error
	GetByID(ctx context.Context, id int64) (*domain.Movie, error)
	Update(ctx context.Context, movie *domain.Movie) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context, title string, genres []string) ([]*domain.Movie, error)
}

type MovieService interface {
	CreateMovie(ctx context.Context, movie *domain.Movie) error
	GetMovieByID(ctx context.Context, id int64) (*domain.Movie, error)
	GetAllMovie(ctx context.Context, title string, genres []string) ([]*domain.Movie, error)
	UpdateMovie(ctx context.Context, movie *domain.Movie) error
	DeleteMovie(ctx context.Context, id int64) error
}
