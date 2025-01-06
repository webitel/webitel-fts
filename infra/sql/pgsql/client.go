package pgsql

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webitel/webitel-fts/infra/sql"
	"github.com/webitel/wlog"
)

type DB struct {
	ctx  context.Context
	pool *pgxpool.Pool
	log  *wlog.Logger
}

type rows struct {
	pgx.Rows
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
	r, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return &rows{
		Rows: r,
	}, nil
}

func (db *DB) Close() error {
	db.pool.Close()
	return nil
}

func (r *rows) Columns() []string {
	c := make([]string, 0, len(r.FieldDescriptions()))
	for _, v := range r.FieldDescriptions() {
		c = append(c, v.Name)
	}

	return c
}
