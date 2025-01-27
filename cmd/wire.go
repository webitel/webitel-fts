//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	"github.com/google/wire"
	"github.com/webitel/webitel-fts/config"
	"github.com/webitel/webitel-fts/internal/handler"
	"github.com/webitel/webitel-fts/internal/service"
	"github.com/webitel/webitel-fts/internal/store"
)

var wireAppResourceSet = wire.NewSet(
	search, log, grpcSrv, pubsubConn, setupCluster,
)

var wireAppHandlersSet = wire.NewSet(
	store.NewIndexEngine,
	service.NewIndexEngine, wire.Bind(new(service.IndexEngineStore), new(*store.IndexEngine)),
	handler.NewSearchEngine, wire.Bind(new(handler.SearchEngineService), new(*service.IndexEngine)),

	handler.NewSubscriber, wire.Bind(new(handler.SubscriberService), new(*service.IndexEngine)),
)

func initAppResources(context.Context, *config.Config) (*resources, func(), error) {
	wire.Build(wireAppResourceSet, wire.Struct(new(resources),
		"search", "log", "grpcSrv", "pubsub", "cluster"))
	return &resources{}, nil, nil
}

func initAppHandlers(*resources) (*handlers, error) {
	wire.Build(wireAppHandlersSet,
		wire.FieldsOf(new(*resources), "search", "log", "grpcSrv", "pubsub"),
		wire.Struct(new(handlers), "searchEngine", "subscriber"),
	)

	return &handlers{}, nil
}

// Setup cmd
var wireSetupResourceSet = wire.NewSet(
	search, log, pubsubConn,
)

var wireSetupHandlersSet = wire.NewSet(
	store.NewManagement,
	service.NewManagement, wire.Bind(new(service.ManagementStore), new(*store.Management)),
	handler.NewManagement, wire.Bind(new(handler.ManagementService), new(*service.Management)),
)

func initSetupResources(context.Context, *config.Config) (*resources, func(), error) {
	wire.Build(wireSetupResourceSet, wire.Struct(new(resources),
		"log", "search"))
	return &resources{}, nil, nil
}

func initSetupHandlers(*resources) (*handlers, error) {
	wire.Build(wireSetupHandlersSet,
		wire.FieldsOf(new(*resources), "search", "log"),
		wire.Struct(new(handlers), "management"),
	)

	return &handlers{}, nil
}

// Import cmd
var wireImportResourceSet = wire.NewSet(
	search, log, setupSql,
)

var wireImportHandlersSet = wire.NewSet(
	handler.NewImport,
)

func initImportResources(context.Context, *config.Config) (*resources, func(), error) {
	wire.Build(wireImportResourceSet, wire.Struct(new(resources),
		"log", "search", "sql"))
	return &resources{}, nil, nil
}

func initImportHandlers(*resources) (*handlers, error) {
	wire.Build(wireImportHandlersSet,
		wire.FieldsOf(new(*resources), "search", "log", "sql"),
		wire.Struct(new(handlers), "importData"),
	)

	return &handlers{}, nil
}
