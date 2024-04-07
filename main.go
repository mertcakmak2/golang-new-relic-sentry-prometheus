package main

import (
	"context"
	"errors"
	"fmt"
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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger = config.ZapTestConfig()

// @title           Go Monitoring App
// @version         1.0
// @description     Go HTTP server with Gin framework.

// @BasePath  /
func main() {

	// Postgres Config & Migration
	db := config.ConnectPostgres()
	database.Migrate(db)
	defer func() {
		logger.Info("DB connection closing...")
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

	srv := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Server Shutdown: %s\n", err))
	}

	select {
	case <-ctx.Done():
		logger.Info("timeout of 5 seconds.")
	}
	logger.Info("Server exiting")
}

func setupRouter(newRelicConfig *newrelic.Application, handler *user.Handler) *gin.Engine {
	router := gin.Default()

	// Swagger => http://localhost:8080/swagger/index.html
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
