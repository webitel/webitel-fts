package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/webitel/webitel-fts/config"
	"github.com/webitel/wlog"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	cfg *config.Config
	log *wlog.Logger
	ctx context.Context
	eg  errgroup.Group
}

func NewApp(cfg *config.Config, ctx context.Context) *App {
	return &App{
		cfg: cfg,
		log: wlog.GlobalLogger(),
		ctx: ctx,
	}
}

func (a *App) Run() (func(), error) {

	r, shutdown, err := initAppResources(a.ctx, a.cfg)
	if err != nil {
		return nil, err
	}

	a.log = r.log

	_, err = initAppHandlers(r)
	if err != nil {
		return shutdown, err
	}

	a.eg.Go(func() error {
		a.log.Debug("listen grpc " + r.grpcSrv.Addr)
		return r.grpcSrv.Listen()
	})

	a.eg.Go(func() error {
		a.log.Debug("start pubsub")
		return r.pubsub.Start()
	})

	return shutdown, nil
}

func apiCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "server",
		Aliases: []string{"a"},
		Usage:   "Start FTS API server",
		Flags:   apiFlags(cfg),
		Action: func(c *cli.Context) error {
			interruptChan := make(chan os.Signal, 1)

			app := NewApp(cfg, c.Context)
			shutdown, err := app.Run()
			defer func() {
				if shutdown != nil {
					shutdown()
				}
			}()
			if err != nil {
				wlog.Error(err.Error(), wlog.Err(err))
				return err
			}
			signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			<-interruptChan
			return nil
		},
	}
}

func apiFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "service-id",
			Category:    "server",
			Usage:       "service id ",
			Value:       "1",
			Destination: &cfg.Service.Id,
			Aliases:     []string{"i"},
			EnvVars:     []string{"ID"},
		},
		&cli.StringFlag{
			Name:        "bind-address",
			Category:    "server",
			Usage:       "address that should be bound to for internal cluster communications",
			Value:       "127.0.0.1:10011",
			Destination: &cfg.Service.Address,
			Aliases:     []string{"b"},
			EnvVars:     []string{"BIND_ADDRESS"},
		},
		&cli.StringFlag{
			Name:        "consul-discovery",
			Category:    "server",
			Usage:       "service discovery address",
			Value:       "127.0.0.1:8500",
			Destination: &cfg.Service.Consul,
			Aliases:     []string{"c"},
			EnvVars:     []string{"CONSUL"},
		},
		&cli.StringFlag{
			Name:        "pubsub",
			Category:    "service/pubsub",
			Usage:       "publish/subscribe rabbitmq broker connection string",
			Value:       "amqp://webitel:webitel@127.0.0.1:5672/",
			Destination: &cfg.Pubsub.Address,
			Aliases:     []string{"p"},
			EnvVars:     []string{"PUBSUB"},
		},
	}
}
