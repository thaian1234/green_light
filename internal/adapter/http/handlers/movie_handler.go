package handlers

import (
	"github.com/gin-gonic/gin"
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
	movie, err := h.movieSvc.GetMovieByID(req.ID)
	if err != nil {
		response.SendNotFound(ctx, err)
	}
	response.SendSuccess(ctx, response.Envelope{
		"movie": movie,
	})
}
