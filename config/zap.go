package config

import (
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func ZapConfig() *zap.Logger {
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zap.InfoLevel)

	backgroundCore, err := nrzap.WrapBackgroundCore(core, app)
	if err != nil && err != nrzap.ErrNilApp {
		panic(err)
	}

	return zap.New(backgroundCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
