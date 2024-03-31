package config

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
)

func NewRelicConfig() *newrelic.Application {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(config().NewRelic.AppName),
		newrelic.ConfigLicense(config().NewRelic.License),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if nil != err {
		fmt.Printf("New Relic initialization failed: %v", err)
		os.Exit(1)
	}

	return app
}
