package delivery

import (
	customerpb "DobrikaDev/customer-service/internal/generated/proto/customer"
	"DobrikaDev/customer-service/internal/service/customer"
	"DobrikaDev/customer-service/utils/config"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	customerService *customer.CustomerService
	customerpb.UnimplementedCustomerServiceServer

	cfg    *config.Config
	logger *zap.Logger
}

func NewServer(ctx context.Context, customerService *customer.CustomerService, cfg *config.Config, logger *zap.Logger) *Server {
	server := &Server{customerService: customerService, cfg: cfg, logger: logger}
	return server
}

func (s *Server) Register(grpcServer *grpc.Server) {
	customerpb.RegisterCustomerServiceServer(grpcServer, s)
}
