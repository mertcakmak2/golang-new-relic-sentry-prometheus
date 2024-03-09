package main

import (
	"errors"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"go-app/config"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	// Sentry Config
	config.SentryConfig()
	// New Relic Config
	app := config.NewRelicConfig()
	// Zap Config
	logger = config.ZapConfig()

	router := gin.Default()

	router.Use(nrgin.Middleware(app))
	router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	router.Use(func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.Scope().SetTag("someRandomTag", "maybeYouNeedIt")
			hub.Scope().SetTag("requestId", "123-asdf")
		}
		ctx.Next()
	})

	v1 := router.Group("/v1")
	v1.GET("/login", v1login)

	router.Run(":8000")
}

func v1login(c *gin.Context) {
	err := errors.New("failure exception")
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("unwantedQuery", "someQueryDataMaybe")
			hub.CaptureMessage("User provided unwanted query string, but we recovered just fine")
			hub.CaptureException(err)
		})
	}
	logger.Info("v1 login request received", zap.String("test", "test"))

	logger.Error(err.Error(), zap.Error(err))
	c.Writer.WriteString("v1 login")
}
