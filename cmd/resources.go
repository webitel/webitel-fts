package cmd

import (
	"context"
	"github.com/webitel/webitel-fts/config"
	"github.com/webitel/webitel-fts/infra/grpc"
	"github.com/webitel/webitel-fts/infra/pubsub"
	"github.com/webitel/webitel-fts/infra/searchengine"
	"github.com/webitel/webitel-fts/infra/searchengine/opensearch"
	"github.com/webitel/webitel-fts/infra/sql"
	"github.com/webitel/webitel-fts/infra/sql/pgsql"
	"github.com/webitel/webitel-fts/internal/handler"
	"github.com/webitel/wlog"
	"strings"
)

type handlers struct {
	searchEngine *handler.SearchEngine
	subscriber   *handler.Subscriber

	management *handler.Management
	importData *handler.ImportData
}

type resources struct {
	search  searchengine.SearchEngine
	log     *wlog.Logger
	grpcSrv *grpc.Server
	pubsub  *pubsub.Manager
	sql     sql.Store
}

func grpcSrv(cfg *config.Config, l *wlog.Logger) (*grpc.Server, func(), error) {
	s, err := grpc.New(cfg.Service.Address, l)
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

	l := wlog.NewLogger(logConfig)
	wlog.InitGlobalLogger(l)

	return l, func() {
		// TODO sync otel
	}, nil
}

func setupSql(log *wlog.Logger, cfg *config.Config) (sql.Store, func(), error) {
	s, err := pgsql.New(context.Background(), cfg.SqlSettings.DSN, log)
	if err != nil {
		return nil, nil, err
	}

	return s, func() {
		s.Close()
	}, nil
}
