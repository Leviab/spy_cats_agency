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
	"context"
	"log/slog"
	"net/http"
	"os"

	"spy_cats_agency/internal/config"
	"spy_cats_agency/internal/handler"
	"spy_cats_agency/internal/repository/postgres"
	"spy_cats_agency/internal/service"
	"spy_cats_agency/pkg/catapi"
	"spy_cats_agency/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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

	// Run migrations
	if err := db.RunMigrations(context.Background()); err != nil {
		appLogger.Error("Failed to run migrations", slog.Any("error", err))
		os.Exit(1)
	}
	appLogger.Info("Database migrations completed")

	// Initialize the CatAPI client
	catAPIClient := catapi.NewClient(cfg.CatAPIEndpoint)

	// Initialize services
	catService := service.NewCatService(db, catAPIClient)
	missionService := service.NewMissionService(db, db, db) // db implements all repo interfaces

	// Initialize handlers
	catHandler := handler.NewCatHandler(catService)
	missionHandler := handler.NewMissionHandler(missionService)

	// Set up router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.Middleware(appLogger))
	router.Use(handler.ErrorMiddleware())

	// API v1 routes
	api := router.Group("/api/v1")
	{
		// Cat routes
		cats := api.Group("/cats")
		{
			cats.POST("", catHandler.CreateCat)
			cats.GET("", catHandler.ListCats)
			cats.GET("/:id", catHandler.GetCat)
			cats.PATCH("/:id/salary", catHandler.UpdateCatSalary)
			cats.DELETE("/:id", catHandler.DeleteCat)
		}

		// Mission routes
		missions := api.Group("/missions")
		{
			missions.POST("", missionHandler.CreateMission)
			missions.GET("", missionHandler.ListMissions)
			missions.GET("/:id", missionHandler.GetMission)
			missions.DELETE("/:id", missionHandler.DeleteMission)
			missions.PATCH("/:id/assign-cat", missionHandler.AssignCatToMission)
			missions.PATCH("/:id/complete", missionHandler.CompleteMission)
			missions.POST("/:id/targets", missionHandler.AddTargetToMission)
		}

		// Target routes
		targets := api.Group("/targets")
		{
			targets.PATCH("/:id/notes", missionHandler.UpdateTargetNotes)
			targets.PATCH("/:id/complete", missionHandler.CompleteTarget)
			targets.DELETE("/:id", missionHandler.DeleteTarget)
		}

	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	serverAddr := ":" + cfg.ServerPort
	appLogger.Info("Server starting", slog.String("address", serverAddr))

	// Start the server
	if err := router.Run(serverAddr); err != nil {
		appLogger.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
