package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/webitel/webitel-fts/config"
	"os"
)

type ImportConfig struct {
}

func importCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "import",
		Aliases: []string{"i"},
		Usage:   "Import data",
		Flags:   importFlags(cfg),
		Action: func(c *cli.Context) error {
			return importData(cfg)
		},
	}
}

func importFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "index",
			Category:    "import",
			Usage:       "Import to index",
			Value:       "",
			Destination: &cfg.Import.Index,
			Aliases:     []string{"in"},
		},
		&cli.StringFlag{
			Name:        "query",
			Category:    "import",
			Usage:       "Import query to index",
			Value:       "",
			Destination: &cfg.Import.Query,
			Aliases:     []string{"q"},
		},
		&cli.StringFlag{
			Name:        "col-domain",
			Category:    "import",
			Usage:       "domain id from query",
			Value:       "dc",
			Destination: &cfg.Import.ColDomainId,
			Aliases:     []string{"d"},
		},
		&cli.StringFlag{
			Name:        "col-id",
			Category:    "import",
			Usage:       "document id from query",
			Value:       "id",
			Destination: &cfg.Import.ColId,
			Aliases:     []string{"id"},
		},
		&cli.StringFlag{
			Name:        "file",
			Category:    "import",
			Usage:       "query file",
			Value:       "",
			Destination: &cfg.Import.File,
			Aliases:     []string{"f"},
		},
	}
}

func importData(cfg *config.Config) error {
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
	q := cfg.Import.Query

	if cfg.Import.File != "" {
		data, err := os.ReadFile(cfg.Import.File)
		if err != nil {
			return err
		}
		q = string(data)
	}

	return h.importData.Import(ctx, q, cfg.Import.ColId, cfg.Import.ColDomainId, cfg.Import.Index)
}
