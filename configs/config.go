package configs

import (
	"fmt"
	"time"

	cfg "github.com/Yalantis/go-config"
)

type (
	Config struct {
		APIHost         string        `envconfig:"API_LISTEN_URL"       default:"0.0.0.0:3000"`
		DebugHost       string        `envconfig:"API_DEBUG_URL"        default:"0.0.0.0:4000"`
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
		DisableTLS   bool          `envconfig:"API_POSTGRES_DISABLE_TLS"   default:"true"`
	}
)

// FromFile return config from file path.
func FromFile(path string) (Config, error) {
	var config Config
	err := cfg.Init(&config, path)
	if err != nil {
		return config, fmt.Errorf("init path %s: %w", path, err)
	}

	return config, nil
}
