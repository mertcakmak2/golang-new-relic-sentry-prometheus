package main

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-app/config"
	"go-app/database"
	"go-app/docs"
	"go-app/middleware"
	"go-app/user"
)

var logger = config.ZapTestConfig()

// @title           Go Monitoring App
// @version         1.0
// @description     Go HTTP server with Gin framework.

// @host      localhost:8080
// @BasePath  /
func main() {

	// Postgres Config & Migration
	db := config.ConnectPostgres()
	database.Migrate(db)
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	// Sentry Config, New Relic Config & Zap Config
	config.SentryConfig()
	newRelicConfig := config.NewRelicConfig()
	logger = config.ZapConfig(newRelicConfig)

	// User Repository, User UseCase & User Handler
	userRepo := user.NewUserRepository(db)
	userUseCase := user.NewUserUseCase(userRepo, logger)
	userHandler := user.NewUserHandler(userUseCase, logger)

	// Setup Router
	router := setupRouter(newRelicConfig, userHandler)
	router.Run(":8080")
}

func setupRouter(newRelicConfig *newrelic.Application, handler *user.Handler) *gin.Engine {
	router := gin.Default()

	// Swagger => http://localhost:8080/swagger/index.html
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	// Middlewares
	_middleware := middleware.NewMiddleware(newRelicConfig, logger)
	router.Use(_middleware.NewRelicMiddleWare())
	router.Use(_middleware.SentryMiddleware())
	router.Use(_middleware.LogMiddleware)

	// Prometheus Metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Endpoints
	v1 := router.Group("/api/v1/users")
	v1.POST("", handler.CreateUser)
	v1.GET("/:id", handler.GetUserById)
	v1.PUT("", handler.UpdateUser)
	v1.DELETE("/:id", handler.DeleteUserById)

	return router
}
