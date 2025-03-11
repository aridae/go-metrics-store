package grpc

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/pkg/logger"
	metricsstorepb "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
	"google.golang.org/grpc"
	"net"
)

type grpcServer interface {
	grpc.ServiceRegistrar
	Stop()
	Serve(net.Listener) error
}

type Server struct {
	server grpcServer
	port   int
}

func NewServer(port int, apiServer metricsstorepb.MetricsStoreAPIServer) *Server {
	grpcServ := grpc.NewServer()
	metricsstorepb.RegisterMetricsStoreAPIServer(grpcServ, apiServer)

	return &Server{
		server: grpcServ,
		port:   port,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.server.Stop()
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to start tcp listener <port:%d>: %w", s.port, err)
	}
	logger.Infof("start listening on address %v", listener.Addr())

	if err = s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	return nil
}
