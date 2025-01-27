package grpc

import (
	"context"
	"errors"
	"github.com/webitel/webitel-fts/infra/webitel"
	"github.com/webitel/webitel-fts/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const hdrTokenAccess = "X-Webitel-Access"

type grpcSessionKey struct {
}

func authUnaryInterceptor(api *webitel.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		token, err := tokenFromContext(ctx)
		if err != nil {
			return nil, err
		}

		session, err := api.CachedSession(ctx, token)
		if err != nil {
			return nil, err
		}

		return handler(context.WithValue(ctx, grpcSessionKey{}, session), req)
	}
}

func tokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("empty metadata")
	}

	token := md.Get(hdrTokenAccess)
	if len(token) < 1 {
		return "", errors.New("can't find authorization token")
	}

	if token[0] == "" {
		return "", errors.New("empty authorization token")
	}

	return token[0], nil
}

func SessionFromContext(ctx context.Context) model.Session {
	session, ok := ctx.Value(grpcSessionKey{}).(model.Session)
	if !ok {
		return model.Session{}
	}

	return session
}
