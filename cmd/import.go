package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/webitel/webitel-fts/config"
)

type ImportConfig struct {
}

func importCmd(cfg *config.Config) *cli.Command {
	var index string
	return &cli.Command{
		Name:    "import",
		Aliases: []string{"i"},
		Usage:   "Import data",
		Flags:   importFlags(cfg, &index),
		Action: func(c *cli.Context) error {
			q := `select id, domain_id, name, created_by, updated_by, created_at, updated_at, calendar, dnc_list,
       team, description, schema, task_processing, tags
from call_center.cc_queue_list;`
			colId := "id"
			colDom := "domain_id"
			return importData(cfg, q, colId, colDom, index)
		},
	}
}

func importFlags(cfg *config.Config, index *string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "index",
			Category:    "storage",
			Usage:       "Import to index",
			Value:       "",
			Destination: index,
			Aliases:     []string{"in"},
		},
	}
}

func importData(cfg *config.Config, q string, colId string, colDomainId string, index string) error {
	ctx := context.Background()
	r, cl, err := initImportResources(ctx, cfg)
	if err != nil {
		return err
	}

	defer cl()

	h, err := initImportHandlers(r)
	if err != nil {
		return err
	}

	return h.importData.Import(ctx, q, colId, colDomainId, index)
}
