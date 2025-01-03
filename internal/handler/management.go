package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/wlog"
	"io"
	"os"
)

type ManagementService interface {
	UpsertTemplate(ctx context.Context, t *model.Template) error
}

type Management struct {
	svc ManagementService
}

func NewManagement(svc ManagementService) *Management {
	h := &Management{
		svc: svc,
	}
	return h
}

func (m *Management) UpsertTemplate(ctx context.Context, template string) error {
	f, err := os.Open(template)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var tms []*model.Template
	err = json.Unmarshal(data, &tms)
	if err != nil {
		return err
	}

	for _, t := range tms {
		err = m.svc.UpsertTemplate(ctx, t)
		if err != nil {
			return err
		}
		wlog.Info(fmt.Sprintf("upsert template %s - success", t.Name))
	}

	return nil
}
