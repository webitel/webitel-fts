package service

import (
	"context"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/wlog"
)

type SearchEngineStore interface {
	Search(ctx context.Context, user *model.SignedInUser, in *model.SearchQuery) ([]*model.SearchResult, error)
}

type SearchEngine struct {
	store SearchEngineStore
	log   *wlog.Logger
}

func NewSearchEngine(store SearchEngineStore, log *wlog.Logger) *SearchEngine {
	return &SearchEngine{
		store: store,
		log:   log.With(),
	}
}

func (s *SearchEngine) Search(ctx context.Context, user *model.SignedInUser, search *model.SearchQuery) ([]*model.SearchResult, bool, error) {
	search.Limit++
	next := false
	res, err := s.store.Search(ctx, user, search)
	if err != nil {
		return nil, false, err
	}

	if len(res) > search.Limit-1 {
		next = true
		res = res[:search.Limit-1]
	}

	return res, next, nil
}
