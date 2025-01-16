package http

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/config"
	"github.com/thaian1234/green_light/internal/adapter/http/handlers"
)

type Routes struct {
	r   *gin.Engine
	cfg *config.Config
}

func NewRoutes(
	r *gin.Engine,
	cfg *config.Config,
	healthHandler *handlers.HealthHandler,
	movieHandler *handlers.MovieHandler,
) (*Routes, error) {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ginConfig := cors.DefaultConfig()
	allowedOrigins := cfg.HTTP.AllowedOrigins
	originsList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	r.Use(cors.New(ginConfig))
	r.NoRoute(gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	}))
	r.NoMethod(gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Method not allowed"})
	}))

	v1 := r.Group("/v1/api")
	{
		// Health route
		health := v1.Group("/health")
		{
			health.GET("/", healthHandler.Check)
		}
		// Movie route
		movie := v1.Group("/movie/")
		{
			movie.GET("/:id", movieHandler.ShowMovie)
			movie.GET("/", movieHandler.ListMovies)
			movie.POST("/", movieHandler.CreateMovie)
			movie.PATCH("/:id", movieHandler.UpdateMovie)
			movie.DELETE("/:id", movieHandler.DeleteMovie)
		}
	}

	return &Routes{
		r:   r,
		cfg: cfg,
	}, nil
}
