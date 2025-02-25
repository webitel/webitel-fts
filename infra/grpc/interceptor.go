package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/webitel/webitel-fts/infra/webitel"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/wlog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
	"time"
)

const hdrTokenAccess = "X-Webitel-Access"
const requestContextName = "grpc_ctx"

type grpcSessionKey struct {
}

type ApiErr struct {
	Id            string
	Code          int
	Status        string `json:"status"` // Message to be display to the end user without debugging information
	DetailedError string `json:"detail"` // Internal error string to help the developer
}

func internalErr(err error) string {
	a := ApiErr{
		Id:            "fts.app.err",
		Code:          http.StatusInternalServerError,
		Status:        "",
		DetailedError: err.Error(),
	}
	s, _ := json.Marshal(a)
	return string(s)
}

func authErr(err error) string {
	a := ApiErr{
		Id:            "fts.auth.err",
		Code:          http.StatusForbidden,
		Status:        "",
		DetailedError: err.Error(),
	}
	s, _ := json.Marshal(a)
	return string(s)
}

func authUnaryInterceptor(logger *wlog.Logger, api *webitel.Client) grpc.UnaryServerInterceptor {
	tc := NewTrace()

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		var session model.Session
		var res any

		start := time.Now()
		md, token, err := tokenFromContext(ctx)
		if md == nil {
			md = metadata.MD{}
		}
		ip := getClientIp(md)

		if err != nil {
			logger.Error(err.Error(),
				wlog.Err(err),
				wlog.String("ip_address", ip),
			)
			return nil, status.Error(codes.PermissionDenied, authErr(err))
		}

		propagators := otel.GetTextMapPropagator()
		ctx = propagators.Extract(
			ctx, GrpcHeaderCarrier(md),
		)

		spanCtx, span := tc.Start(ctx, info.FullMethod)
		defer func() {
			span.End()
		}()

		session, err = api.CachedSession(spanCtx, token)
		span.SetAttributes(
			attribute.Int64("domain_id", session.DomainId),
			attribute.Int64("user_id", session.UserId),
			attribute.String("ip_address", ip),
			attribute.String("method", info.FullMethod),
		)

		if err != nil {
			span.SetStatus(otelCodes.Error, err.Error())
			logger.Error(err.Error(),
				wlog.Err(err),
				wlog.String("ip_address", ip),
			)
			return nil, status.Error(codes.PermissionDenied, authErr(err))
		}

		log := logger.With(wlog.Namespace("context"),
			wlog.Int64("domain_id", session.DomainId),
			wlog.Int64("user_id", session.UserId),
			wlog.String("ip_address", ip),
			wlog.String("method", info.FullMethod),
		)

		res, err = handler(context.WithValue(ctx, grpcSessionKey{}, session), req)
		if err != nil {
			span.SetStatus(otelCodes.Error, err.Error())
			log.Error(err.Error(), wlog.Float64("duration_ms", float64(time.Since(start).Microseconds())/float64(1000)))
			return nil, status.Error(codes.Internal, internalErr(err))
		} else {
			span.SetStatus(otelCodes.Ok, "success")
			log.Debug(strings.Join(md.Get("url"), ","),
				wlog.Float64("duration_ms", float64(time.Since(start).Microseconds())/float64(1000)))
		}

		return res, nil
	}
}

func tokenFromContext(ctx context.Context) (metadata.MD, string, error) {
	var info metadata.MD
	var ok bool
	v := ctx.Value(requestContextName)
	info, ok = v.(metadata.MD)

	if !ok {
		info, ok = metadata.FromIncomingContext(ctx)
	}

	if !ok {
		return info, "", errors.New("empty metadata")
	}

	token := info.Get(hdrTokenAccess)
	if len(token) < 1 {
		return info, "", errors.New("can't find authorization token")
	}

	if token[0] == "" {
		return info, "", errors.New("empty authorization token")
	}

	return info, token[0], nil
}

func SessionFromContext(ctx context.Context) model.Session {
	session, ok := ctx.Value(grpcSessionKey{}).(model.Session)
	if !ok {
		return model.Session{}
	}

	return session
}

func getClientIp(info metadata.MD) string {
	ip := strings.Join(info.Get("x-real-ip"), ",")
	if ip == "" {
		ip = strings.Join(info.Get("x-forwarded-for"), ",")
	}

	return ip
}
