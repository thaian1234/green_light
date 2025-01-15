package services

import (
	"context"

	"github.com/thaian1234/green_light/internal/core/domain"
	"github.com/thaian1234/green_light/internal/core/ports"
)

type MovieService struct {
	movieRepo ports.MovieRepository
}

func NewMovieService(movieRepo ports.MovieRepository) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
	}
}

func (s *MovieService) CreateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	err := s.movieRepo.Insert(ctx, movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) GetMovieByID(ctx context.Context, id int64) (*domain.Movie, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) GetAllMovie(ctx context.Context) ([]*domain.Movie, error) {
	return s.movieRepo.GetAll(ctx)
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	err := s.movieRepo.Update(ctx, movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, id int64) error {
	return s.movieRepo.Delete(ctx, id)
}
