package config

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"os"
)

func SentryConfig() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              config().Sentry.Dsn,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v", err)
		os.Exit(1)
	}
}
