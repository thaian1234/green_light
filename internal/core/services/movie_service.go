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

func (s *MovieService) CreateMovie(movie *domain.Movie) (*domain.Movie, error) {
	return &domain.Movie{}, nil
}

func (s *MovieService) GetMovieByID(id int64) (*domain.Movie, error) {
	return &domain.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "My movie",
		Year:      2024,
		Genres:    []string{"Horror", "Action"},
		Runtime:   1,
		Version:   1,
	}, nil
}
