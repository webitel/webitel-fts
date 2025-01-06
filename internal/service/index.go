package service

import (
	"context"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/wlog"
)

type IndexEngineStore interface {
	Create(ctx context.Context, msg model.Message) error
	Update(ctx context.Context, msg model.Message) error
	Delete(ctx context.Context, msg model.Message) error
	Search(ctx context.Context, domainId int64, in *model.SearchQuery) ([]*model.SearchResult, error)
}

type IndexEngine struct {
	store IndexEngineStore
	log   *wlog.Logger
}

func NewIndexEngine(log *wlog.Logger, s IndexEngineStore) *IndexEngine {
	return &IndexEngine{
		log: log.With(
			wlog.String("service", "index"),
		),
		store: s,
	}
}

func (s *IndexEngine) Create(ctx context.Context, msg model.Message) error {
	return s.store.Create(ctx, msg)
}

func (s *IndexEngine) Update(ctx context.Context, msg model.Message) error {
	return s.store.Update(ctx, msg)
}

func (s *IndexEngine) Delete(ctx context.Context, msg model.Message) error {
	return s.store.Delete(ctx, msg)
}

func (s *IndexEngine) Search(ctx context.Context, user *model.SignedInUser, search *model.SearchQuery) ([]*model.SearchResult, bool, error) {
	search.Limit++
	next := false
	res, err := s.store.Search(ctx, user.DomainId, search)
	if err != nil {
		return nil, false, err
	}

	if len(res) > search.Limit-1 {
		next = true
		res = res[:search.Limit-1]
	}

	return res, next, nil
}
