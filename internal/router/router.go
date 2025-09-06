package router

import (
	"log/slog"
	"spy_cats_agency/internal/handler"
	"spy_cats_agency/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Config holds the dependencies needed for route setup.
type Config struct {
	CatHandler     *handler.CatHandler
	MissionHandler *handler.MissionHandler
	TargetHandler  *handler.TargetHandler
	Logger         *slog.Logger
}

// Setup initializes and configures all routes.
func Setup(cfg Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.Middleware(cfg.Logger))
	router.Use(handler.ErrorMiddleware(cfg.Logger))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// API v1 routes
	setupAPIRoutes(router, cfg)

	return router
}

// setupAPIRoutes configures all API v1 routes.
func setupAPIRoutes(router *gin.Engine, cfg Config) {
	api := router.Group("/api/v1")
	{
		setupCatRoutes(api, cfg.CatHandler)
		setupMissionRoutes(api, cfg.MissionHandler, cfg.TargetHandler)
		setupTargetRoutes(api, cfg.TargetHandler)
	}
}

// setupCatRoutes configures cat-related routes.
func setupCatRoutes(api *gin.RouterGroup, catHandler *handler.CatHandler) {
	cats := api.Group("/cats")
	{
		cats.POST("", catHandler.CreateCat)
		cats.GET("", catHandler.ListCats)
		cats.GET("/:id", catHandler.GetCat)
		cats.PATCH("/:id/salary", catHandler.UpdateCatSalary)
		cats.DELETE("/:id", catHandler.DeleteCat)
	}
}

// setupMissionRoutes configures mission-related routes.
func setupMissionRoutes(api *gin.RouterGroup, missionHandler *handler.MissionHandler, targetHandler *handler.TargetHandler) {
	missions := api.Group("/missions")
	{
		missions.POST("", missionHandler.CreateMission)
		missions.GET("", missionHandler.ListMissions)
		missions.GET("/:id", missionHandler.GetMission)
		missions.DELETE("/:id", missionHandler.DeleteMission)
		missions.PATCH("/:id/assign-cat", missionHandler.AssignCatToMission)
		missions.PATCH("/:id/complete", missionHandler.CompleteMission)
		missions.POST("/:id/targets", targetHandler.AddTargetToMission)
	}
}

// setupTargetRoutes configures target-related routes.
func setupTargetRoutes(api *gin.RouterGroup, targetHandler *handler.TargetHandler) {
	targets := api.Group("/targets")
	{
		targets.PATCH("/:id/notes", targetHandler.UpdateTargetNotes)
		targets.PATCH("/:id/complete", targetHandler.CompleteTarget)
		targets.DELETE("/:id", targetHandler.DeleteTarget)
	}
}
