package service

import (
	"context"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/wlog"
)

type ManagementStore interface {
	UpsertTemplate(ctx context.Context, t *model.Template) error
}

type Management struct {
	store ManagementStore
	log   *wlog.Logger
}

func NewManagement(store ManagementStore, log *wlog.Logger) *Management {
	return &Management{
		store: store,
		log:   log.With(),
	}
}

func (m *Management) UpsertTemplate(ctx context.Context, t *model.Template) error {
	return m.store.UpsertTemplate(ctx, t)
}
