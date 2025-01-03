package pgsql

import (
	"github.com/urfave/cli/v2"
	"github.com/webitel/webitel-fts/config"
)

func Flags(cfg *config.Config) []cli.Flag {
	const category = "database/sql"
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "postgresql-dsn",
			Category:    category,
			Usage:       "Postgres connection string",
			EnvVars:     []string{"DATA_SOURCE"},
			Value:       "postgres://postgres:postgres@localhost:5432/webitel?sslmode=disable",
			Destination: &cfg.SqlSettings.DSN,
		},
	}
}
