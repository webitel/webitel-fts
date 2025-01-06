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
			q := `select c.id, c.dc,
       c.about, c.given_name, c.middle_name, c.family_name,
       c.common_name,
       (select json_agg(p.number) phones
        from contacts.contact_phone p
        where p.contact_id = c.id),
       (select json_agg(e.email) email
        from contacts.contact_email e
        where e.contact_id = c.id)
from contacts.contact c`

			colId := "id"
			colDom := "dc"
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
