package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yourusername/projectname/internal/client"
	"github.com/yourusername/projectname/internal/config"
	"github.com/yourusername/projectname/internal/handler"
	"github.com/yourusername/projectname/internal/middleware"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/internal/service"
	"github.com/yourusername/projectname/pkg/logger"
	"gorm.io/gorm"
)

// Router wraps gin.Engine and provides additional functionality
type Router struct {
	*gin.Engine
	db     *gorm.DB
	logger *logger.Logger
	config *config.Config
}

// NewRouter creates and configures a new router instance with all required dependencies
func NewRouter(cfg *config.Config, logger *logger.Logger, authClient client.AuthClient, db *gorm.DB) *Router {
	gin.SetMode(gin.ReleaseMode)
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// Create CSRF store
	csrfStore := middleware.NewCSRFStore()

	// Add base middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware(logger))

	router := &Router{
		Engine: r,
		db:     db,
		logger: logger,
		config: cfg,
	}

	router.setupRoutes(authClient, csrfStore)
	return router
}

// setupRoutes configures all the routes for the application
func (r *Router) setupRoutes(authClient client.AuthClient, csrfStore *middleware.CSRFStore) {
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// Web routes (with CSRF protection)
	web := r.Group("")
	web.Use(middleware.CSRFMiddleware(csrfStore))
	{
		// Add web routes here that need CSRF protection
	}

	// API routes (no CSRF protection)
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Public routes (no auth required)
			public := v1.Group("")
			{
				public.GET("/health", r.HealthCheck)

				// Community routes
				communityRepo := repository.NewCommunityRepository(r.db)
				communityService := service.NewCommunityService(communityRepo)
				communityHandler := handler.NewCommunityHandler(communityService)
				communityHandler.RegisterRoutes(public)

			}

			// Protected routes (require authentication)
			protected := v1.Group("")
			protected.Use(
				middleware.AuthMiddleware(authClient, r.logger),
				middleware.TenancyMiddleware(),
			)
			{
				// User routes
				userRepo := repository.NewUserRepository(r.db)
				userService := service.NewUserService(userRepo)
				userHandler := handler.NewUserHandler(userService)
				userHandler.RegisterRoutes(protected)

				// Community routes
				//communityRepo := repository.NewCommunityRepository(r.db)
				// := service.NewCommunityService(communityRepo)
				//communityHandler := handler.NewCommunityHandler(communityService)
				//communityHandler.RegisterRoutes(protected)

				// Building routes
				//buildingRepo := repository.NewBuildingRepository(r.db)
				//buildingService := service.NewBuildingService(buildingRepo)
				//buildingHandler := handler.NewBuildingHandler(buildingService)
				//buildingHandler.RegisterRoutes(protected)

				// Unit routes
				//unitRepo := repository.NewUnitRepository(r.db)
				//unitService := service.NewUnitService(unitRepo)
				//unitHandler := handler.NewUnitHandler(unitService)
				//unitHandler.RegisterRoutes(protected)

				// Admin routes (require admin role)
				admin := protected.Group("/admin")
				admin.Use(middleware.RequireRoles("ADMIN"))
				{
					// Add admin routes here
				}
			}
		}
	}
}

// HealthCheck handles the health check endpoint
// @Summary Health check endpoint
// @Description Get the health status of the API
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (r *Router) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
