package service

import (
	"context"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/wlog"
)

type IndexEngineStore interface {
	Create(ctx context.Context, msg model.Message) error
}

type Subscriber struct {
	store IndexEngineStore
	log   *wlog.Logger
}

func NewSubscriber(log *wlog.Logger, s IndexEngineStore) *Subscriber {
	return &Subscriber{
		log:   log,
		store: s,
	}
}

func (s *Subscriber) Create(ctx context.Context, msg model.Message) error {
	return s.store.Create(ctx, msg)
}
