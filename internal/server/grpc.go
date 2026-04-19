package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/grpc-healthd/internal/health"
)

type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	checker *health.Checker
}

func NewHealthServer(checker *health.Checker) *HealthServer {
	return &HealthServer{checker: checker}
}

func (s *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	service := req.GetService()
	status, err := s.checker.GetStatus(service)
	if err != nil {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN,
		}, nil
	}
	var grpcStatus grpc_health_v1.HealthCheckResponse_ServingStatus
	if status.Healthy {
		grpcStatus = grpc_health_v1.HealthCheckResponse_SERVING
	} else {
		grpcStatus = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	}
	return &grpc_health_v1.HealthCheckResponse{Status: grpcStatus}, nil
}

func (s *HealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return grpc.ErrServerStopped
}

func ListenAndServe(addr string, hs *HealthServer) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, hs)
	return grpcServer.Serve(lis)
}
