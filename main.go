package main

import (
	"DobrikaDev/customer-service/di"
	customerpb "DobrikaDev/customer-service/internal/generated/proto/customer"
	"DobrikaDev/customer-service/utils/config"
	"DobrikaDev/customer-service/utils/logger"
	"context"
	"os"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoadConfigFromFile("deployments/config.yaml")
	logger, _ := logger.NewLogger()
	defer logger.Sync()

	container := di.NewContainer(ctx, cfg, logger)

	customerpb.RegisterCustomerServiceServer(
		container.GetGRPCServer(),
		container.GetRpcServer(),
	)

	logger.Info("Starting application with port", zap.String("port", cfg.Port))

	err := container.GetGRPCServer().Serve(*container.GetNetListener())
	if err != nil {
		logger.Error("Error while serving grpcServer:", zap.Error(err))
		os.Exit(1)
	}
}
