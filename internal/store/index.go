package store

import (
	"context"
	"fmt"
	"github.com/webitel/webitel-fts/infra/searchengine"
	"github.com/webitel/webitel-fts/internal/model"
)

type IndexEngine struct {
	db searchengine.SearchEngine
}

func NewIndexEngine(d searchengine.SearchEngine) *IndexEngine {
	return &IndexEngine{
		db: d,
	}
}

func (s *IndexEngine) Create(ctx context.Context, msg model.Message) error {
	return s.db.Insert(ctx,
		fmt.Sprintf("%v", msg.Id),
		fmt.Sprintf("%v_%v", msg.ObjectName, msg.DomainId),
		msg.Body,
	)
}
