package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contains environment variables for the application, database, cache, token, logger and http server
type (
	Config struct {
		App     *App
		Token   *Token
		Redis   *Redis
		DB      *DB
		HTTP    *HTTP
		Logger  *Logger
		Limiter *Limiter
	}
	// App contains all the environment variables for the application
	App struct {
		Name    string
		Env     string
		Version string
	}
	// Token contains all the environment variables for the token service
	Token struct {
		Duration string
	}
	// Redis contains all the environment variables for the cache service
	Redis struct {
		Addr     string
		Password string
	}
	// Database contains all the environment variables for the database
	DB struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
	}
	// HTTP contains all the environment variables for the http server
	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
	}
	// Logger configuration
	Logger struct {
		LogPath     string
		LogLevel    string
		LogMaxSize  int
		LogBackUps  int
		LogMaxAge   int
		LogCompress bool
	}
	// Limiter configuration
	Limiter struct {
		Rps     int
		Burst   int
		Enabled bool
	}
)

// Load creates a new container instance
func Load() (*Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name:    os.Getenv("APP_NAME"),
		Env:     os.Getenv("APP_ENV"),
		Version: os.Getenv("APP_VERSION"),
	}

	token := &Token{
		Duration: os.Getenv("TOKEN_DURATION"),
	}

	redis := &Redis{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	db := &DB{
		Connection: os.Getenv("DB_CONNECTION"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Name:       os.Getenv("DB_NAME"),
	}

	http := &HTTP{
		Env:            os.Getenv("APP_ENV"),
		URL:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	logMaxSize, _ := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
	logBackUps, _ := strconv.Atoi(os.Getenv("LOG_BACKUPS"))
	logMaxAge, _ := strconv.Atoi(os.Getenv("LOG_MAX_AGE"))
	logCompress, _ := strconv.ParseBool(os.Getenv("LOG_COMPRESS"))

	logger := &Logger{
		LogPath:     os.Getenv("LOG_PATH"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
		LogMaxSize:  logMaxSize,
		LogBackUps:  logBackUps,
		LogMaxAge:   logMaxAge,
		LogCompress: logCompress,
	}

	rps, _ := strconv.Atoi(os.Getenv("LIMITER_RPS"))
	burst, _ := strconv.Atoi(os.Getenv("LIMITER_BURST"))
	enabled, _ := strconv.ParseBool(os.Getenv("LIMITER_ENABLED"))
	limiter := &Limiter{
		Rps:     rps,
		Burst:   burst,
		Enabled: enabled,
	}

	return &Config{
		app,
		token,
		redis,
		db,
		http,
		logger,
		limiter,
	}, nil
}
