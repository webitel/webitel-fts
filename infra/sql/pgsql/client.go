package pgsql

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webitel/webitel-fts/infra/sql"
	"github.com/webitel/wlog"
)

type DB struct {
	ctx  context.Context
	pool *pgxpool.Pool
	log  *wlog.Logger
}

func New(ctx context.Context, dsn string, log *wlog.Logger) (sql.Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	db := &DB{
		ctx:  ctx,
		pool: pool,
		log:  log,
	}

	return db, nil
}

func (db *DB) Select(ctx context.Context, out any, query string, args ...any) error {
	return pgxscan.Select(ctx, db.pool, out, query, args...)
}

func (db *DB) Query(ctx context.Context, sql string, args ...any) (sql.Rows, error) {
	return db.pool.Query(ctx, sql, args...)
}

func (db *DB) Close() error {
	db.pool.Close()
	return nil
}
