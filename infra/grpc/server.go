package grpc

import (
	"github.com/webitel/wlog"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	Addr string
	log  *wlog.Logger
	*grpc.Server
}

// New provides a new gRPC server.
func New(addr string, log *wlog.Logger) (*Server, error) {

	s := grpc.NewServer()

	return &Server{
		Addr:   addr,
		Server: s,
		log:    log,
	}, nil
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func (s *Server) Shutdown() error {
	s.log.Debug("receive shutdown grpc ")
	s.Server.GracefulStop()

	return nil
}
