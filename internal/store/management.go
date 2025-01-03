package store

import (
	"context"
	"github.com/webitel/webitel-fts/infra/searchengine"
	"github.com/webitel/webitel-fts/internal/model"
)

type Management struct {
	db searchengine.SearchEngine
}

func NewManagement(d searchengine.SearchEngine) *Management {
	return &Management{
		db: d,
	}
}

func (s *Management) UpsertTemplate(ctx context.Context, t *model.Template) error {
	return s.db.Template(ctx, t.Name, t.Data)
}
