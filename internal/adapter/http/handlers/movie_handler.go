package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/internal/core/domain"
	"github.com/thaian1234/green_light/internal/core/ports"
)

type MovieHandler struct {
	movieSvc ports.MovieService
}

func NewMovieHandler(movieSvc ports.MovieService) *MovieHandler {
	return &MovieHandler{
		movieSvc: movieSvc,
	}
}

type (
	params struct {
		ID int64 `uri:"id" binding:"required,min=1,number"`
	}
	createMovieRequest struct {
		Title   string   `json:"title" binding:"required"`
		Year    int32    `json:"year" binding:"required,year_range"`
		Runtime int32    `json:"runtime" binding:"required,number,min=1"`
		Genres  []string `json:"genres" binding:"required"`
	}
	updateMovieRequest struct {
		Title   *string  `json:"title" binding:"omitempty,required"`
		Year    *int32   `json:"year" binding:"omitempty,required,year_range"`
		Runtime *int32   `json:"runtime" binding:"omitempty,required"`
		Genres  []string `json:"genres" binding:"omitempty,required"`
	}
)

func (h *MovieHandler) ShowMovie(ctx *gin.Context) {
	var req params
	if err := ctx.ShouldBindUri(&req); err != nil {
		HandleValidationError(ctx, err)
		return
	}
	movie, err := h.movieSvc.GetMovieByID(ctx, req.ID)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	SendSuccess(ctx, Envelope{
		"movie": movie,
	})
}

func (h *MovieHandler) CreateMovie(ctx *gin.Context) {
	var req createMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleValidationError(ctx, err)
		return
	}
	movieModal := domain.Movie{
		Title:   req.Title,
		Year:    req.Year,
		Runtime: domain.Runtime(req.Runtime),
		Genres:  req.Genres,
	}
	createdMovie, err := h.movieSvc.CreateMovie(ctx, &movieModal)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	SendCreatedSuccess(ctx, createdMovie)
}

func (h *MovieHandler) ListMovies(ctx *gin.Context) {
	movies, err := h.movieSvc.GetAllMovie(ctx)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	SendSuccess(ctx, Envelope{
		"movies": movies,
	})
}

func (h *MovieHandler) UpdateMovie(ctx *gin.Context) {
	var param params
	if err := ctx.ShouldBindUri(&param); err != nil {
		HandleValidationError(ctx, err)
		return
	}

	var reqBody updateMovieRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		HandleValidationError(ctx, err)
		return
	}
	var movieModel domain.Movie
	movieModel.ID = param.ID

	if reqBody.Title != nil {
		movieModel.Title = *reqBody.Title
	}
	if reqBody.Year != nil {
		movieModel.Year = *reqBody.Year
	}
	if reqBody.Runtime != nil {
		movieModel.Runtime = domain.Runtime(*reqBody.Runtime)
	}
	if reqBody.Genres != nil {
		movieModel.Genres = reqBody.Genres
	}

	updatedMovie, err := h.movieSvc.UpdateMovie(ctx, &movieModel)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	SendUpdatedSuccess(ctx, Envelope{
		"movie": updatedMovie,
	})
}

func (h *MovieHandler) DeleteMovie(ctx *gin.Context) {
	var req params
	if err := ctx.ShouldBindUri(&req); err != nil {
		HandleValidationError(ctx, err)
		return
	}
	err := h.movieSvc.DeleteMovie(ctx, req.ID)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	SendDeletedSuccess(ctx)
}
