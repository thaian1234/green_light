package services

import (
	"time"

	"github.com/thaian1234/green_light/internal/core/domain"
)

type MovieService struct {
}

func NewMovieService() *MovieService {
	return &MovieService{}
}

func (s *MovieService) CreateMovie(movie domain.Movie) domain.Movie {
	return domain.Movie{}
}

func (s *MovieService) GetMovieByID(id int64) domain.Movie {
	return domain.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "My movie",
		Year:      2024,
		Genres:    []string{"Horror", "Action"},
		Runtime:   102,
		Version:   1,
	}
}
