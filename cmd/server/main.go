package main

import (
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/nsvirk/moneybotsapi/internal/api"
	"github.com/nsvirk/moneybotsapi/internal/config"
	"github.com/nsvirk/moneybotsapi/internal/repository"
	"github.com/nsvirk/moneybotsapi/internal/utils/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Get()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Print the configuration
	fmt.Printf("%+v", cfg)

	// Connect to Postgres
	db, err := repository.ConnectPostgres(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Postgres: %v", err))
	}

	// Initialize logger
	if err := logger.InitLogger(db); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	// Log application startup
	log.Info().
		Str("module", "main").
		Str("version", cfg.APIVersion).
		Str("environment", cfg.ServerEnv).
		Str("extra_key", "extra_value").
		Msg(cfg.APIName + " started")

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Setup routes
	api.SetupRoutes(e, cfg, db)

	// Log server start
	log.Info().
		Str("module", "main").
		Str("port", cfg.ServerPort).
		Msg("Starting HTTP server")

	// Start server
	if err := startServer(e, cfg); err != nil {
		log.Error().
			Str("module", "main").
			Err(err).
			Msg("Server failed to start")
		os.Exit(1)
	}
}

func startServer(e *echo.Echo, cfg *config.Config) error {
	port := cfg.ServerPort
	if port == "" {
		port = "3007"
		log.Warn().
			Str("module", "main").
			Msg("No port configured, using default: 3007")
	}

	// Log all incoming requests
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			log.Info().
				Str("module", "http").
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Str("remote_ip", c.RealIP()).
				Msg("API request")

			err := next(c)

			// Log error if it occurred
			if err != nil {
				log.Error().
					Str("module", "http").
					Str("method", c.Request().Method).
					Str("path", c.Request().URL.Path).
					Err(err).
					Dur("latency", time.Since(start)).
					Msg("Request error")
			}

			return err
		}
	})

	return e.Start(":" + port)
}
