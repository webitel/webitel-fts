package grpc

import (
	"github.com/webitel/webitel-fts/infra/webitel"
	"github.com/webitel/wlog"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type Server struct {
	Addr string
	host string
	port int
	log  *wlog.Logger
	*grpc.Server
	listener net.Listener
}

// New provides a new gRPC server.
func New(addr string, log *wlog.Logger, api *webitel.Client) (*Server, error) {

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		authUnaryInterceptor(api),
	))

	h, p, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	port, _ := strconv.Atoi(p)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		Addr:     addr,
		Server:   s,
		log:      log,
		host:     h,
		port:     port,
		listener: l,
	}, nil
}

func (s *Server) Listen() error {
	return s.Serve(s.listener)
}

func (s *Server) Shutdown() error {
	s.log.Debug("receive shutdown grpc")
	err := s.listener.Close()
	s.Server.GracefulStop()
	return err
}

func (s *Server) Host() string {
	return s.host
}

func (s *Server) Port() int {
	return s.port
}
