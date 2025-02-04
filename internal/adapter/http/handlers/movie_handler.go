package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/internal/core/domain"
	"github.com/thaian1234/green_light/internal/core/ports"
	"github.com/thaian1234/green_light/pkg/util"
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
	listMovieRequest struct {
		Title  string `form:"title"`
		Genres string `form:"genres"`
		domain.Filter
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
	err := h.movieSvc.CreateMovie(ctx, &movieModal)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	SendCreatedSuccess(ctx, movieModal)
}

func (h *MovieHandler) ListMovies(ctx *gin.Context) {
	var queryParams listMovieRequest
	queryParams.SortSafeList = []string{"id", "-id", "title", "-title", "year", "-year", "runtime", "-runtime", "genres", "-genres"}
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		HandleValidationError(ctx, err)
		return
	}

	genres := util.ReadCSV(queryParams.Genres, []string{})
	title := ctx.DefaultQuery("title", "")
	filter := domain.Filter{
		Page: util.ReadInt(queryParams.Page, 1),
		Size: util.ReadInt(queryParams.Size, 20),
		Sort: ctx.DefaultQuery("sort", "id"),
	}

	movies, err := h.movieSvc.GetAllMovie(ctx, title, genres, filter)
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

	existingMovie, err := h.movieSvc.GetMovieByID(ctx, param.ID)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	if reqBody.Title != nil {
		existingMovie.Title = *reqBody.Title
	}
	if reqBody.Year != nil {
		existingMovie.Year = *reqBody.Year
	}
	if reqBody.Runtime != nil {
		existingMovie.Runtime = domain.Runtime(*reqBody.Runtime)
	}
	if reqBody.Genres != nil {
		existingMovie.Genres = reqBody.Genres
	}

	err = h.movieSvc.UpdateMovie(ctx, existingMovie)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	SendUpdatedSuccess(ctx, Envelope{
		"movie": existingMovie,
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
