package middleware

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go-app/logging"
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

func (m middleware) LogMiddleware(ctx *gin.Context) {
	var responseBody = logging.HandleResponseBody(ctx.Writer)
	var requestBody = logging.HandleRequestBody(ctx.Request)
	requestId := uuid.NewString()

	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
		hub.Scope().SetTag("requestId", requestId)
		ctx.Writer = responseBody
	}

	ctx.Next()

	logMessage := logging.FormatRequestAndResponse(ctx.Writer, ctx.Request, responseBody.Body.String(), requestId, requestBody)

	if logMessage != "" {
		if isSuccessStatusCode(ctx.Writer.Status()) {
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
