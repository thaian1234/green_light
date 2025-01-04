package ports

import "github.com/thaian1234/green_light/internal/core/domain"

type MovieService interface {
	CreateMovie(movie domain.Movie) (domain.Movie, error)
	GetMovieByID(id int64) (domain.Movie, error)
}
