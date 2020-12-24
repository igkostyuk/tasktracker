package postgres

import (
	"database/sql"
	"net/url"

	"github.com/igkostyuk/tasktracker/app/config"

	// The database driver in use.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
)

func Open(cfg config.Postgres) (*sql.DB, error) {
	q := make(url.Values)
	if cfg.DisableTLS {
		q.Set("sslmode", "disable")
	}
	q.Set("timezone", "utc")

	// nolint:exhaustivestruct
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sql.Open("pgx", u.String())
	if err != nil {
		return nil, errors.Wrap(err, "creating db")
	}
	err = db.Ping()
	if err != nil {
		db.Close()

		return nil, errors.Wrapf(err, "ping db %s", u.String())
	}

	return db, nil
}
