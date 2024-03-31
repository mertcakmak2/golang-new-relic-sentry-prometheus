package main

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-app/config"
	"go-app/database"
	"go-app/middleware"
	"go-app/user"
)

var logger = config.ZapTestConfig()

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

	// Middlewares
	_middleware := middleware.NewMiddleware(newRelicConfig, logger)
	router.Use(_middleware.NewRelicMiddleWare())
	router.Use(_middleware.SentryMiddleware())
	router.Use(_middleware.LogMiddleware)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	v1 := router.Group("/api/v1/users")
	v1.POST("", handler.CreateUser)
	v1.GET("/:id", handler.GetUserById)
	v1.PUT("", handler.UpdateUser)
	v1.DELETE("/:id", handler.DeleteUserById)
	return router
}
