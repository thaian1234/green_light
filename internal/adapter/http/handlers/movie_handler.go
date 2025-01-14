package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/internal/core/domain"
	"github.com/thaian1234/green_light/internal/core/ports"
	"github.com/thaian1234/green_light/pkg/response"
)

type MovieHandler struct {
	movieSvc ports.MovieService
}

func NewMovieHandler(movieSvc ports.MovieService) *MovieHandler {
	return &MovieHandler{
		movieSvc: movieSvc,
	}
}

type getMovieRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (h *MovieHandler) ShowMovie(ctx *gin.Context) {
	var req getMovieRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.SendValidationError(ctx, err)
		return
	}
	movie, err := h.movieSvc.GetMovieByID(ctx, req.ID)
	if err != nil {
		response.SendNotFound(ctx, err)
		return
	}
	response.SendSuccess(ctx, response.Envelope{
		"movie": movie,
	})
}

type createMovieRequest struct {
	Title   string   `json:"title" binding:"required"`
	Year    int32    `json:"year" binding:"required,year_range"`
	Runtime int32    `json:"runtime" binding:"required,number,min=1"`
	Genres  []string `json:"genres" binding:"required"`
}

func (h *MovieHandler) CreateMovie(ctx *gin.Context) {
	var req createMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendValidationError(ctx, err)
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
		response.SendBadRequest(ctx, err)
		return
	}
	response.SendCreatedSuccess(ctx, createdMovie)
}
