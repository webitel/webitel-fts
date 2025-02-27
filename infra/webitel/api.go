package webitel

import (
	"context"
	"fmt"
	"github.com/hashicorp/golang-lru/v2/expirable"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/webitel/webitel-fts/gen/api"
	"github.com/webitel/webitel-fts/internal/model"
	"github.com/webitel/wlog"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

var (
	sessionGroupRequest singleflight.Group
)

const tokenRequestTimeout = time.Second * 15

type Client struct {
	sessionCache *expirable.LRU[string, *model.Session]
	auth         api.AuthClient
}

func NewClient(consulTarget string, log *wlog.Logger) (*Client, error) {
	authConn, err := grpc.Dial(fmt.Sprintf("consul://%s/go.webitel.app?wait=14s", consulTarget),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		auth:         api.NewAuthClient(authConn),
		sessionCache: expirable.NewLRU[string, *model.Session](10000, nil, time.Second*10), // TODO config
	}, nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) CachedSession(ctx context.Context, token string) (model.Session, error) {
	s, ok := c.sessionCache.Get(token)
	if ok {
		return *s, nil
	}

	result, err, shared := sessionGroupRequest.Do(token, func() (interface{}, error) {
		ctx2, cancel := context.WithTimeout(ctx, tokenRequestTimeout)
		defer cancel()

		return c.GetSession(ctx2, token)
	})

	if err != nil {
		return model.Session{}, err
	}

	s = result.(*model.Session)
	if !shared {
		c.sessionCache.Add(token, s)
	}

	return *s, nil
}

func (c *Client) GetSession(ctx context.Context, token string) (*model.Session, error) {
	header := metadata.New(map[string]string{"x-webitel-access": token})
	authCtx := metadata.NewOutgoingContext(ctx, header)
	res, err := c.auth.UserInfo(authCtx, nil)
	if err != nil {
		return nil, err
	}

	s := &model.Session{
		Name:     res.Username,
		DomainId: res.Dc,
		Expire:   res.ExpiresAt,
		UserId:   res.UserId,
		Scopes:   nil,
	}

	hasAdminRead := false
	for _, v := range res.Permissions {
		if v.Id == "read" {
			hasAdminRead = true
			break
		}
	}

	for _, scope := range res.Scope {
		s.Scopes = append(s.Scopes, model.SessionPermission{
			Id:     scope.Id,
			Class:  scope.Class,
			Obac:   scope.Obac && !hasAdminRead,
			Rbac:   scope.Rbac && !hasAdminRead,
			Access: scope.Access,
		})
	}

	for _, role := range res.Roles {
		s.RoleIds = append(s.RoleIds, role.Id)
	}
	s.RoleIds = append(s.RoleIds, res.UserId)

	return s, nil
}
