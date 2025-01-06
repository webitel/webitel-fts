package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/webitel/webitel-fts/config"
	"github.com/webitel/wlog"
)

func setupCmd(cfg *config.Config) *cli.Command {
	var template string
	return &cli.Command{
		Name:    "setup",
		Aliases: []string{"s"},
		Usage:   "Setup FTS templates",
		Flags:   setupFlags(cfg, &template),
		Action: func(c *cli.Context) error {
			return setupTemplate(c.Context, cfg, template)
		},
	}
}

func setupFlags(cfg *config.Config, template *string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "template",
			Category:    "file",
			Usage:       "JSON index template",
			Value:       "./mapping/contacts.json",
			Destination: template,
			Aliases:     []string{"t"},
			EnvVars:     []string{"INDEX_TEMPLATE"},
		},
	}
}

func setupTemplate(ctx context.Context, cfg *config.Config, template string) error {
	r, shutdown, err := initSetupResources(ctx, cfg)
	if err != nil {
		return err
	}
	defer shutdown()

	h, err := initSetupHandlers(r)
	if err != nil {
		return err
	}

	err = h.management.UpsertTemplate(ctx, template)
	if err != nil {
		wlog.Error(err.Error())
		return err
	}

	return nil
}
