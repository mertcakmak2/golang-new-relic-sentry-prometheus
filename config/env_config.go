package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"log"
	"sync"
)

var (
	cfg        AppConfig
	configOnce sync.Once
)

func config() AppConfig {
	configOnce.Do(func() {
		ctx := context.Background()
		if err := envconfig.Process(ctx, &cfg); err != nil {
			log.Fatal(err)
		}
		log.Println("Environments initialized.")
	})
	return cfg
}

type AppConfig struct {
	Database *Database
	NewRelic *NewRelic
	Sentry   *Sentry
}

type Database struct {
	Host         string `env:"POSTGRES_HOST, default=localhost"`
	Username     string `env:"POSTGRES_USERNAME, default=postgres"`
	Password     string `env:"POSTGRES_PASSWORD, default=postgres"`
	Port         string `env:"POSTGRES_PORT, default=5432"`
	DatabaseName string `env:"DATABASE_NAME, default=postgres"`
}

type NewRelic struct {
	AppName string `env:"APP_NAME, default=go-app"`
	License string `env:"NEW_RELIC_LICENSE"`
}

type Sentry struct {
	Dsn string `env:"SENTRY_DSN"`
}
