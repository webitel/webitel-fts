package cmd

import (
	"context"
	"github.com/webitel/webitel-fts/config"
	"github.com/webitel/webitel-fts/infra/consul"
	"github.com/webitel/webitel-fts/infra/grpc"
	"github.com/webitel/webitel-fts/infra/pubsub"
	"github.com/webitel/webitel-fts/infra/searchengine"
	"github.com/webitel/webitel-fts/infra/searchengine/opensearch"
	"github.com/webitel/webitel-fts/infra/sql"
	"github.com/webitel/webitel-fts/infra/sql/pgsql"
	"github.com/webitel/webitel-fts/infra/webitel"
	"github.com/webitel/webitel-fts/internal/handler"
	"github.com/webitel/webitel-fts/internal/model"
	otelsdk "github.com/webitel/webitel-go-kit/otel/sdk"
	"github.com/webitel/wlog"
	"go.opentelemetry.io/otel/sdk/resource"
	"strings"

	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	// -------------------- plugin(s) -------------------- //
	_ "github.com/webitel/webitel-go-kit/otel/sdk/log/otlp"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/log/stdout"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/metric/otlp"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/metric/stdout"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/trace/otlp"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/trace/stdout"
)

type handlers struct {
	searchEngine *handler.SearchEngine
	subscriber   *handler.Subscriber

	management *handler.Management
	importData *handler.ImportData
}

type resources struct {
	search    searchengine.SearchEngine
	log       *wlog.Logger
	grpcSrv   *grpc.Server
	pubsub    *pubsub.Manager
	sql       sql.Store
	cluster   *consul.Cluster
	apiClient *webitel.Client
}

func grpcSrv(cfg *config.Config, l *wlog.Logger, client *webitel.Client) (*grpc.Server, func(), error) {
	s, err := grpc.New(cfg.Service.Address, l, client)
	if err != nil {
		return nil, nil, err
	}
	return s, func() {
		if err := s.Shutdown(); err != nil {
			l.Error(err.Error(), wlog.Err(err))
		}
	}, nil
}

func search(cfg *config.Config, l *wlog.Logger) (searchengine.SearchEngine, func(), error) {
	s, err := opensearch.New(strings.Split(cfg.OpenSearchConfig.Addresses, ","),
		cfg.OpenSearchConfig.Username,
		cfg.OpenSearchConfig.Password,
		cfg.OpenSearchConfig.Insecure,
	)

	if err != nil {
		return nil, nil, err
	}
	return s, func() {
		if err := s.Shutdown(); err != nil {
			l.Error(err.Error(), wlog.Err(err))
		}
	}, nil

}

func pubsubConn(log *wlog.Logger, cfg *config.Config) (*pubsub.Manager, func(), error) {
	ps, err := pubsub.New(log, cfg.Pubsub.Address)
	if err != nil {
		return nil, nil, err
	}

	return ps, func() {
		if err := ps.Shutdown(); err != nil {
			log.Error(err.Error(), wlog.Err(err))
		}
	}, nil
}

func log(cfg *config.Config) (*wlog.Logger, func(), error) {
	logSettings := cfg.Log

	if !logSettings.Console && !logSettings.Otel && len(logSettings.File) == 0 {
		logSettings.Console = true
	}

	logConfig := &wlog.LoggerConfiguration{
		EnableConsole: logSettings.Console,
		ConsoleJson:   false,
		ConsoleLevel:  logSettings.Lvl,
	}

	if logSettings.File != "" {
		logConfig.FileLocation = logSettings.File
		logConfig.EnableFile = true
		logConfig.FileJson = true
		logConfig.FileLevel = logSettings.Lvl
	}

	var otelShutdownFunc otelsdk.ShutdownFunc
	ctx, cancel := context.WithCancel(context.Background())

	exit := func() {
		if otelShutdownFunc != nil {
			err := otelShutdownFunc(ctx)
			if err != nil {
				wlog.Error(err.Error(), wlog.Err(err))
			}
		}
		cancel()
	}

	if logSettings.Otel {
		// TODO
		var err error
		logConfig.EnableExport = true
		otelShutdownFunc, err = otelsdk.Configure(
			ctx,
			otelsdk.WithResource(resource.NewSchemaless(
				semconv.ServiceName(model.ServiceName),
				semconv.ServiceVersion(model.CurrentVersion),
				semconv.ServiceInstanceID(cfg.Service.Id),
				semconv.ServiceNamespace("webitel"),
			)),
		)
		if err != nil {
			return nil, nil, err
		}
	}

	l := wlog.NewLogger(logConfig)
	wlog.RedirectStdLog(l)
	wlog.InitGlobalLogger(l)

	return l, exit, nil
}

func setupSql(log *wlog.Logger, cfg *config.Config) (sql.Store, func(), error) {
	s, err := pgsql.New(context.Background(), cfg.SqlSettings.DSN, log)
	if err != nil {
		return nil, nil, err
	}

	return s, func() {
		err = s.Close()
		if err != nil {
			wlog.Error(err.Error(), wlog.Err(err))
		}
	}, nil
}

func setupCluster(cfg *config.Config, srv *grpc.Server) (*consul.Cluster, func(), error) {
	c := consul.NewCluster(model.ServiceName, cfg.Service.Consul)
	err := c.Start(cfg.Service.Id, srv.Host(), srv.Port())

	if err != nil {
		return nil, nil, err
	}
	return c, func() {
		c.Stop()
	}, nil

}

func setupApiClient(cfg *config.Config, l *wlog.Logger) (*webitel.Client, func(), error) {
	c, err := webitel.NewClient(cfg.Service.Consul, l)
	if err != nil {
		return nil, nil, err
	}

	return c, func() {
		err = c.Close()
		if err != nil {
			wlog.Error(err.Error(), wlog.Err(err))
		}
	}, nil

}
