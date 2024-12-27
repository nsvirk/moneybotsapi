// Package api contains the API routes for the Moneybots API
package api

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/nsvirk/moneybotsapi/internal/config"
	"github.com/nsvirk/moneybotsapi/pkg/utils/response"
	"gorm.io/gorm"
)

// SetupRoutes configures the routes for the API
func SetupRoutes(e *echo.Echo, cfg *config.Config, db *gorm.DB) {

	// Create a group for all API routes
	api := e.Group("")

	// Index route
	api.GET("/", indexRoute)

}

// indexRoute sets up the index route for the API
func indexRoute(c echo.Context) error {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	message := fmt.Sprintf("%s %s", cfg.APIName, cfg.APIVersion)
	return response.SuccessResponse(c, message)
}
