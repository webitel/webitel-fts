package opensearch

import (
	"github.com/urfave/cli/v2"
	"github.com/webitel/webitel-fts/config"
)

func Flags(cfg *config.Config) []cli.Flag {
	const category = "searchengine/opensearch"

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "opensearch-host",
			Category:    category,
			Usage:       "OpenSearch list of nodes to use.",
			EnvVars:     []string{"OPENSEARCH_ADDR"},
			Value:       "https://localhost:9200",
			Aliases:     []string{"oh"},
			Destination: &cfg.OpenSearchConfig.Addresses,
		},
		&cli.StringFlag{
			Name:        "opensearch-user",
			Category:    category,
			Usage:       "OpenSearch username for HTTP Basic Authentication.",
			EnvVars:     []string{"OPENSEARCH_USER"},
			Value:       "admin",
			Aliases:     []string{"ou"},
			Destination: &cfg.OpenSearchConfig.Username,
		},
		&cli.StringFlag{
			Name:        "opensearch-pass",
			Category:    category,
			Usage:       "OpenSearch password for HTTP Basic Authentication.",
			EnvVars:     []string{"OPENSEARCH_PASS"},
			Value:       "admin",
			Aliases:     []string{"op"},
			Destination: &cfg.OpenSearchConfig.Password,
		},
	}
}
