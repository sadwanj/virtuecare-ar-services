package grpcserver

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// RegisterFunc defines a callback to register service-specific handlers
type RegisterFunc func(s *grpc.Server)

// Config holds server configuration
type Config struct {
	Host string
	Port string
	//DB       *gorm.DB
	Services []string // List of service names for health reporting
}

// Server represents a generic gRPC server
type Server struct {
	server *grpc.Server
	config *Config
	health *health.Server
}

// NewServer creates a new generic gRPC server
func NewServer(config *Config, reg RegisterFunc, opts ...grpc.ServerOption) (*Server, error) {
	// Standard interceptors can be passed via opts or defined here
	grpcServer := grpc.NewServer(opts...)
	healthServer := health.NewServer()

	s := &Server{
		server: grpcServer,
		config: config,
		health: healthServer,
	}

	// Register Health and Reflection (standard for all services)
	grpc_health_v1.RegisterHealthServer(s.server, s.health)
	reflection.Register(s.server)

	// Call the service-specific registration logic
	reg(s.server)

	return s, nil
}

// Start starts the server and sets health status
func (s *Server) Start() error {
	addr := net.JoinHostPort(s.config.Host, s.config.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	//logger.Infof("gRPC server listening on %s", addr)

	// Set all services to SERVING
	s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	for _, svc := range s.config.Services {
		s.health.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_SERVING)
	}

	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.health.Shutdown()
	s.server.GracefulStop()
}

// HealthCheck provides a standard DB ping check for now just retrn nil
func (s *Server) HealthCheck(ctx context.Context) error {
	// if s.config.DB == nil {
	// 	return errors.New("database connection is nil")
	// }

	// sqlDB, err := s.config.DB.DB()
	// if err != nil {
	// 	return err
	// }

	return nil //sqlDB.PingContext(ctx)
}
