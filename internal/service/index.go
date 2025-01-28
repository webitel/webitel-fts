package service

import (
	"context"
	"fmt"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/wlog"
	"slices"
	"strings"
)

type IndexEngineStore interface {
	Create(ctx context.Context, msg model.Message) error
	Update(ctx context.Context, msg model.Message) error
	Delete(ctx context.Context, msg model.Message) error
	Search(ctx context.Context, domainId int64, in *model.SearchQuery) ([]*model.SearchResult, error)
	GetSupportObjectsName() ([]string, error)
}

type IndexEngine struct {
	store   IndexEngineStore
	log     *wlog.Logger
	objects []string
}

func NewIndexEngine(log *wlog.Logger, s IndexEngineStore) *IndexEngine {
	i := &IndexEngine{
		log: log.With(
			wlog.String("service", "index"),
		),
		store: s,
	}
	var err error
	i.objects, err = s.GetSupportObjectsName()
	if err != nil {
		// TODO schedule refresh
		i.log.Error(err.Error(), wlog.Err(err))
	} else {
		i.log.Debug(fmt.Sprintf("handle [%v] objects", strings.Join(i.objects, ", ")))
	}
	return i
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

func (s *IndexEngine) Search(ctx context.Context, session *model.Session, search *model.SearchQuery) ([]*model.SearchResult, bool, error) {
	search.Limit++
	next := false

	allowObjects := Filter(s.objects, func(e string) bool {
		p := session.ObjectPermission(e)
		if p != nil {
			return p.HasRead()
		}
		return false
	})

	if len(search.ObjectsName) == 0 {
		search.ObjectsName = append([]string{}, allowObjects...)
	}

	search.ObjectsName = Filter(search.ObjectsName, func(e string) bool {
		return slices.Contains(allowObjects, e)
	})

	res, err := s.store.Search(ctx, session.DomainId, search)
	if err != nil {
		return nil, false, err
	}

	if len(res) > search.Limit-1 {
		next = true
		res = res[:search.Limit-1]
	}

	return res, next, nil
}

func Filter[E any](s []E, f func(E) bool) []E {
	s2 := make([]E, 0, len(s))
	for _, e := range s {
		if f(e) {
			s2 = append(s2, e)
		}
	}
	return s2
}
