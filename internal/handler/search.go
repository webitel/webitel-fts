package handler

import (
	"context"
	"fmt"
	pb "github.com/webitel/webitel-fts/gen/go/api/fts"
	"github.com/webitel/webitel-fts/infra/grpc"
	"github.com/webitel/webitel-fts/infra/pubsub"
	"github.com/webitel/webitel-fts/internal/model"
)

type SearchEngineService interface {
	Search(ctx context.Context, user *model.SignedInUser, search *model.SearchQuery) ([]*model.SearchResult, bool, error)
}

type SearchEngine struct {
	pb.UnsafeFTSServiceServer
	svc SearchEngineService
}

func NewSearchEngine(svc SearchEngineService, s *grpc.Server, _ *pubsub.Manager) *SearchEngine {
	h := &SearchEngine{
		svc: svc,
	}
	pb.RegisterFTSServiceServer(s, h)
	return h
}

func (h *SearchEngine) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	items, next, err := h.svc.Search(ctx, &model.SignedInUser{
		DomainId: 1,
	}, &model.SearchQuery{
		Q:           in.GetQ(),
		ObjectsName: in.GetObjectName(),
		Limit:       int(in.GetLimit()),
	})

	if err != nil {
		return nil, err
	}

	var data []*pb.SearchData

	for _, v := range items {
		data = append(data, &pb.SearchData{
			Id:         fmt.Sprintf("%v", v.Id),
			ObjectName: v.ObjectName,
			Text:       v.Text,
		})
	}

	return &pb.SearchResponse{
		Next:  next,
		Items: data,
	}, nil
}
