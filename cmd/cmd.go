package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/webitel/webitel-fts/config"
	"github.com/webitel/webitel-fts/infra/searchengine/opensearch"
	"github.com/webitel/webitel-fts/infra/sql/pgsql"
	"github.com/webitel/webitel-fts/internal/model"
	"os"
	"time"
)

var (

	// version is the APP's semantic version.
	version = model.CurrentVersion

	// commit is the git commit used to build the api.
	Commit     = "hash"
	CommitDate = time.Now().UTC().String()

	Branch = "branch"
	Build  = "dev"
)

// Run the default command
func Run() error {
	cfg := &config.Config{}

	def := &cli.App{
		Name:      "webitel-fts",
		Usage:     "FTS in the Webitel",
		Version:   fmt.Sprintf("%s-%s, %s@%s at %s", version, Build, Branch, Commit, CommitDate),
		Compiled:  time.Now(),
		Copyright: "Webitel, 2024",
		Action: func(c *cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			apiCmd(cfg),
			setupCmd(cfg),
			importCmd(cfg),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Category:    "observability/logging",
				Usage:       "application log level",
				Value:       "debug",
				Destination: &cfg.Log.Lvl,
				Aliases:     []string{"l"},
			},
		},
	}

	def.Flags = append(def.Flags, opensearch.Flags(cfg)...)
	def.Flags = append(def.Flags, pgsql.Flags(cfg)...)
	if err := def.Run(os.Args); err != nil {
		return err
	}

	return nil
}
