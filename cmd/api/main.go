// @title Spy Cat Agency API
// @version 1.0
// @description API for managing spy cats, missions, and targets
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"log/slog"
	"os"

	"spy_cats_agency/internal/config"
	"spy_cats_agency/internal/handler"
	"spy_cats_agency/internal/repository/postgres"
	"spy_cats_agency/internal/router"
	"spy_cats_agency/internal/service"
	"spy_cats_agency/pkg/catapi"
	"spy_cats_agency/pkg/logger"

	"github.com/joho/godotenv"

	_ "spy_cats_agency/docs" // This is required for Swagger
)

func main() {
	appLogger := logger.New()

	// Load .env file for local development
	if err := godotenv.Load(); err != nil {
		appLogger.Info("No .env file found, using environment variables")
	}

	cfg, err := config.LoadConfig(".")
	if err != nil {
		appLogger.Error("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	db, err := postgres.New(cfg)
	if err != nil {
		appLogger.Error("Failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	appLogger.Info("Database connection established")

	// Initialize repositories
	catRepo := postgres.NewCatRepository(db)
	missionRepo := postgres.NewMissionRepository(db)
	targetRepo := postgres.NewTargetRepository(db)

	// Initialize the CatAPI client
	catAPIClient := catapi.NewClient(cfg.CatAPIEndpoint)

	// Initialize services
	catService := service.NewCatService(catRepo, catAPIClient)
	missionService := service.NewMissionService(missionRepo, targetRepo, catRepo)

	// Initialize handlers
	catHandler := handler.NewCatHandler(catService)
	missionHandler := handler.NewMissionHandler(missionService)

	// Set up router with all routes
	routerInstance := router.Setup(router.Config{
		CatHandler:     catHandler,
		MissionHandler: missionHandler,
		Logger:         appLogger,
	})

	serverAddr := ":" + cfg.ServerPort
	appLogger.Info("Server starting", slog.String("address", serverAddr))

	// Start the server
	if err := routerInstance.Run(serverAddr); err != nil {
		appLogger.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
