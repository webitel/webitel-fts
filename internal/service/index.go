package service

import (
	"context"
	"fmt"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/webitel-fts/pkg/client"
	"github.com/webitel/wlog"
	"strings"
)

type IndexEngineStore interface {
	Create(ctx context.Context, msg client.Message) error
	Update(ctx context.Context, msg client.Message) error
	Delete(ctx context.Context, msg client.Message) error
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

func (s *IndexEngine) Create(ctx context.Context, msg client.Message) error {
	return s.store.Create(ctx, msg)
}

func (s *IndexEngine) Update(ctx context.Context, msg client.Message) error {
	return s.store.Update(ctx, msg)
}

func (s *IndexEngine) Delete(ctx context.Context, msg client.Message) error {
	return s.store.Delete(ctx, msg)
}

func (s *IndexEngine) Search(ctx context.Context, session *model.Session, search *model.SearchQuery) ([]*model.SearchResult, bool, error) {
	search.Limit++
	next := false

	var allowObjects []model.ObjectName
	sq := len(search.ObjectsName)

	for _, v := range s.objects {
		if sq != 0 && !search.HasObject(v) {
			continue
		}
		p := session.ObjectPermission(v)
		if p != nil && p.HasRead() {
			a := model.ObjectName{
				Name:    v,
				RoleIds: nil,
			}
			if p.Rbac {
				a.RoleIds = session.RoleIds
			}
			allowObjects = append(allowObjects, a)
		}
	}
	search.ObjectsName = allowObjects

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
