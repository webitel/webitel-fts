package handler

import (
	"context"
	"fmt"
	pb "github.com/webitel/webitel-fts/gen/api/fts"
	"github.com/webitel/webitel-fts/infra/grpc"
	"github.com/webitel/webitel-fts/infra/webitel"
	"github.com/webitel/webitel-fts/internal/model"
)

type SearchEngineService interface {
	Search(ctx context.Context, user *model.Session, search *model.SearchQuery) ([]*model.SearchResult, bool, error)
}

type SearchEngine struct {
	pb.UnsafeFTSServiceServer
	svc    SearchEngineService
	client *webitel.Client
}

func NewSearchEngine(svc SearchEngineService, s *grpc.Server, apiCli *webitel.Client) *SearchEngine {
	h := &SearchEngine{
		svc:    svc,
		client: apiCli,
	}
	pb.RegisterFTSServiceServer(s, h)
	return h
}

func (h *SearchEngine) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	session := grpc.SessionFromContext(ctx)

	q := &model.SearchQuery{
		Q:           in.GetQ(),
		ObjectsName: nil,
		Limit:       int(in.GetSize()),
		Page:        int(in.GetPage()),
	}
	for _, v := range in.GetObjectName() {
		q.ObjectsName = append(q.ObjectsName, model.ObjectName{
			Name: v,
		})
	}

	items, next, err := h.svc.Search(ctx, &session, q)

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
