package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/aszaychik/prima-inti-api/internal/auth"
	"github.com/aszaychik/prima-inti-api/internal/config"
	"github.com/aszaychik/prima-inti-api/internal/errors"
	"github.com/aszaychik/prima-inti-api/internal/health"
	"github.com/aszaychik/prima-inti-api/internal/middleware"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(handlers *Handlers, authService auth.Service, cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := gin.New()

	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	skipPaths := config.GetSkipPaths(cfg.App.Environment)
	loggerConfig := middleware.NewLoggerConfig(
		cfg.Logging.GetLogLevel(),
		skipPaths,
	)
	router.Use(middleware.Logger(loggerConfig))
	router.Use(errors.ErrorHandler())
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	router.Use(cors.New(corsConfig))

	var checkers []health.Checker
	if cfg.Health.DatabaseCheckEnabled {
		dbChecker := health.NewDatabaseChecker(db)
		checkers = append(checkers, dbChecker)
	}
	healthService := health.NewService(checkers, cfg.App.Version, cfg.App.Environment)
	healthHandler := health.NewHandler(healthService)

	router.GET("/health", healthHandler.Health)
	router.GET("/health/live", healthHandler.Live)
	router.GET("/health/ready", healthHandler.Ready)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rlCfg := cfg.Ratelimit
	if rlCfg.Enabled {
		router.Use(
			middleware.NewRateLimitMiddleware(
				rlCfg.Window,
				rlCfg.Requests,
				func(c *gin.Context) string {
					ip := c.ClientIP()
					if ip == "" {
						ip = c.GetHeader("X-Forwarded-For")
						if ip == "" {
							ip = c.GetHeader("X-Real-IP")
						}
						if ip == "" {
							ip = "unknown"
						}
					}
					return ip
				},
				nil,
			),
		)
	}

	v1 := router.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", handlers.User.Register)
			authGroup.POST("/login", handlers.User.Login)
			authGroup.POST("/refresh", handlers.User.RefreshToken)
			authGroup.POST("/logout", auth.AuthMiddleware(authService), handlers.User.Logout)
			authGroup.GET("/me", auth.AuthMiddleware(authService), handlers.User.GetMe)
		}

		// User endpoints - authenticated users can access their own resources
		usersGroup := v1.Group("/users")
		usersGroup.Use(auth.AuthMiddleware(authService))
		{
			usersGroup.GET("/:id", handlers.User.GetUser)
			usersGroup.PUT("/:id", handlers.User.UpdateUser)
			usersGroup.DELETE("/:id", handlers.User.DeleteUser)
		}

		// Admin endpoints - admin role required, following REST best practices
		adminGroup := v1.Group("/admin")
		adminGroup.Use(auth.AuthMiddleware(authService), middleware.RequireAdmin())
		{
			// User management endpoints
			adminGroup.GET("/users", handlers.User.ListUsers)
			adminGroup.GET("/users/:id", handlers.User.GetUser)
			adminGroup.PUT("/users/:id", handlers.User.UpdateUser)
			adminGroup.DELETE("/users/:id", handlers.User.DeleteUser)
		}

		// ---- Company Profile Routes ----
		v1.GET("/company-profile", handlers.Company.GetCompanyProfile)

		companyGroup := v1.Group("/company-profile")
		companyGroupAdmin := companyGroup.Use(auth.AuthMiddleware(authService), middleware.RequireAdmin())
		{
			companyGroupAdmin.POST("", handlers.Company.CreateCompanyProfile)
			companyGroupAdmin.PUT("", handlers.Company.UpdateCompanyProfile)
			companyGroupAdmin.POST("/links", handlers.Company.AddExternalLink)
			companyGroupAdmin.PUT("/links/:linkId", handlers.Company.UpdateExternalLink)
			companyGroupAdmin.DELETE("/links/:linkId", handlers.Company.DeleteExternalLink)
		}

		// ---- Category Routes ----
		v1.GET("/categories", handlers.Category.ListCategories)
		v1.GET("/categories/:id", handlers.Category.GetCategory)
		v1.GET("/categories/:id/products", handlers.Product.ListProductsByCategory)

		categoryGroup := v1.Group("/categories")
		categoryGroupAdmin := categoryGroup.Use(auth.AuthMiddleware(authService), middleware.RequireAdmin())
		{
			categoryGroupAdmin.POST("", handlers.Category.CreateCategory)
			categoryGroupAdmin.PUT("/:id", handlers.Category.UpdateCategory)
			categoryGroupAdmin.DELETE("/:id", handlers.Category.DeleteCategory)
		}

		// ---- Brand Routes (using consistent :brandId param) ----
		// More specific routes first
		v1.GET("/brands/:brandId/series", handlers.Series.ListSeriesByBrand)
		v1.GET("/brands/:brandId/products", handlers.Product.ListProductsByBrand)
		v1.GET("/brands", handlers.Brand.ListBrands)
		v1.GET("/brands/:brandId", handlers.Brand.GetBrand)

		brandGroup := v1.Group("/brands")
		brandGroupAdmin := brandGroup.Use(auth.AuthMiddleware(authService), middleware.RequireAdmin())
		{
			brandGroupAdmin.POST("", handlers.Brand.CreateBrand)
			brandGroupAdmin.PUT("/:brandId", handlers.Brand.UpdateBrand)
			brandGroupAdmin.DELETE("/:brandId", handlers.Brand.DeleteBrand)
		}

		// ---- Series Routes ----
		v1.GET("/series", handlers.Series.ListSeries)
		v1.GET("/series/:id", handlers.Series.GetSeries)
		v1.GET("/series/:id/products", handlers.Product.ListProductsBySeries)

		seriesAdmin := v1.Group("/series")
		seriesAdmin.Use(auth.AuthMiddleware(authService), middleware.RequireAdmin())
		{
			seriesAdmin.POST("", handlers.Series.CreateSeries)
			seriesAdmin.PUT("/:id", handlers.Series.UpdateSeries)
			seriesAdmin.DELETE("/:id", handlers.Series.DeleteSeries)
		}

		// ---- Product Routes (generic) ----
		v1.GET("/products", handlers.Product.ListProducts)
		v1.GET("/products/:id", handlers.Product.GetProduct)

		productAdmin := v1.Group("/products")
		productAdmin.Use(auth.AuthMiddleware(authService), middleware.RequireAdmin())
		{
			productAdmin.POST("", handlers.Product.CreateProduct)
			productAdmin.PUT("/:id", handlers.Product.UpdateProduct)
			productAdmin.DELETE("/:id", handlers.Product.DeleteProduct)
		}
	}

	return router
}
