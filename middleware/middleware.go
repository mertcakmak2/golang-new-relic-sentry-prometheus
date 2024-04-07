package middleware

import (
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prometheus/client_golang/prometheus"
	"go-app/logging"
	"go-app/metrics"
	"go.uber.org/zap"
	"net/http"
)

type middleware struct {
	newRelicConfig *newrelic.Application
	logger         *zap.Logger
}

func NewMiddleware(newRelicConfig *newrelic.Application, logger *zap.Logger) middleware {
	return middleware{newRelicConfig: newRelicConfig, logger: logger}
}

func (m middleware) NewRelicMiddleWare() gin.HandlerFunc {
	return nrgin.Middleware(m.newRelicConfig)
}

func (m middleware) SentryMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{Repanic: true})
}

/*
Log all HTTP requests and responses to New Relic.
Generates a custom count metric for Prometheus. It uses an HTTP request path and an HTTP method.
*/
func (m middleware) LogMiddleware(ctx *gin.Context) {
	reqMethodAndPath := fmt.Sprintf("[%s] %s", ctx.Request.Method, ctx.FullPath())

	// HTTP Request Response Duration
	timer := prometheus.NewTimer(metrics.HttpRequestDuration.WithLabelValues(reqMethodAndPath))
	defer timer.ObserveDuration()

	var responseBody = logging.HandleResponseBody(ctx.Writer)
	var requestBody = logging.HandleRequestBody(ctx.Request)
	requestId := uuid.NewString()

	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
		hub.Scope().SetTag("requestId", requestId)
		ctx.Writer = responseBody
	}
	ctx.Next()

	statusCode := ctx.Writer.Status()
	// Same HTTP Request Path Counter
	metrics.HttpRequestCountWithPath.With(prometheus.Labels{"url": reqMethodAndPath}).Inc()
	logMessage := logging.FormatRequestAndResponse(statusCode, ctx.Request, responseBody.Body.String(), requestId, requestBody)

	if logMessage != "" {
		if isSuccessStatusCode(statusCode) {
			m.logger.Info(logMessage)
		} else {
			m.logger.Error(logMessage)
		}
	}
}

func isSuccessStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return true
	default:
		return false
	}
}
