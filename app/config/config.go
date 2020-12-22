package config

import (
	"time"

	cfg "github.com/Yalantis/go-config"
	"github.com/pkg/errors"
)

type (
	Config struct {
		APIHost         string        `envconfig:"API_LISTEN_URL"       default:"0.0.0.0:3000"`
		DebugHost       string        `envconfig:"API_LISTEN_URL"       default:"0.0.0.0:4000"`
		ReadTimeout     time.Duration `envconfig:"API_READ_TIMEOUT"     default:"5s"`
		WriteTimeout    time.Duration `envconfig:"API_WRITE_TIMEOUT"    default:"5s"`
		ShutdownTimeout time.Duration `envconfig:"API_SHUTDOWN_TIMEOUT" default:"5s"`

		Postgres Postgres
	}
	Postgres struct {
		Host         string        `envconfig:"POSTGRES_HOST"              default:"0.0.0.0:5432"`
		Name         string        `envconfig:"API_POSTGRES_DATABASE"      default:"postgres"`
		User         string        `envconfig:"API_POSTGRES_USER"          default:"postgres"`
		Password     string        `envconfig:"API_POSTGRES_PASSWORD"      default:"mysecretpassword"`
		PoolSize     int           `envconfig:"API_POSTGRES_POOL_SIZE"     default:"10"`
		MaxRetries   int           `envconfig:"API_POSTGRES_MAX_RETRIES"   default:"5"`
		ReadTimeout  time.Duration `envconfig:"API_POSTGRES_READ_TIMEOUT"  default:"10s"`
		WriteTimeout time.Duration `envconfig:"API_POSTGRES_WRITE_TIMEOUT" default:"10s"`
	}
)

// FromFile return config from file path.
func FromFile(path string) (Config, error) {
	var config Config
	err := cfg.Init(&config, path)
	if err != nil {
		return config, errors.Wrapf(err, "init path: %s", path)
	}

	return config, nil
}