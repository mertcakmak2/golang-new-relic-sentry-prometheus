package config

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
)

var app *newrelic.Application

func NewRelicConfig() *newrelic.Application {
	var err error
	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName("go-strapi-app"),
		newrelic.ConfigLicense("<NEW_RELIC_LICENSE_KEY>"),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if nil != err {
		fmt.Printf("New Relic initialization failed: %v", err)
		os.Exit(1)
	}

	return app
}
