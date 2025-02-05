package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/webitel/webitel-fts/infra/searchengine"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/wlog"
)

type IndexEngine struct {
	log *wlog.Logger
	db  searchengine.SearchEngine
}

func NewIndexEngine(d searchengine.SearchEngine, log *wlog.Logger) *IndexEngine {
	return &IndexEngine{
		db:  d,
		log: log.With(wlog.String("scope", "index_store")),
	}
}

func (s *IndexEngine) Search(ctx context.Context, domainId int64, in *model.SearchQuery) ([]*model.SearchResult, error) {
	var idx []searchengine.IndexSettings

	did := fmt.Sprintf("_%d", domainId)

	for _, v := range in.ObjectsName {
		idx = append(idx, searchengine.IndexSettings{
			Name:          v.Name + did,
			AccessRoleIds: v.RoleIds,
		})
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
			Id:         v.Id,
			ObjectName: v.Index[:len(v.Index)-len(did)], // TODO
			Text:       v.Text,
		})
	}

	return res, nil
}

func (s *IndexEngine) Create(ctx context.Context, msg model.Message) error {
	return s.db.Insert(ctx,
		fmt.Sprintf("%v", msg.Id),
		fmt.Sprintf("%v_%v", msg.ObjectName, msg.DomainId),
		msg.Body,
	)
}

func (s *IndexEngine) Update(ctx context.Context, msg model.Message) error {
	return s.db.Update(ctx,
		fmt.Sprintf("%v", msg.Id),
		fmt.Sprintf("%v_%v", msg.ObjectName, msg.DomainId),
		msg.Body,
	)
}

func (s *IndexEngine) Delete(ctx context.Context, msg model.Message) error {
	return s.db.Delete(ctx,
		fmt.Sprintf("%v", msg.Id),
		fmt.Sprintf("%v_%v", msg.ObjectName, msg.DomainId),
	)
}

func (s *IndexEngine) GetSupportObjectsName() ([]string, error) {
	return s.db.GetTemplates(context.TODO())
}
