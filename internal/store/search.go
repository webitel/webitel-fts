package store

import (
	"context"
	"errors"
	"github.com/webitel/webitel-fts/infra/searchengine"
	"github.com/webitel/webitel-fts/internal/model"
)

type SearchEngine struct {
	db searchengine.SearchEngine
}

func NewSearchEngine(d searchengine.SearchEngine) *SearchEngine {
	return &SearchEngine{
		db: d,
	}
}

func (s *SearchEngine) Search(ctx context.Context, user *model.SignedInUser, in *model.SearchQuery) ([]*model.SearchResult, error) {
	var idx []string
	for _, v := range in.Scope {
		idx = append(idx, v+"_*") // todo add domain
	}

	if len(idx) == 0 {
		return nil, errors.New("scope is required")
	}

	result, err := s.db.Search(ctx, idx, in.Q, in.Limit)
	if err != nil {
		return nil, err
	}

	var res []*model.SearchResult
	for _, v := range result {
		res = append(res, &model.SearchResult{
			Id:    v.Id,
			Scope: v.Index,
			Text:  v.Text,
		})
	}

	return res, nil
}
