package main

import (
	"fmt"

	"log"
	"net"
	"os"

	"github.com/a-dev-mobile/smtp-server/internal/config"
	"github.com/a-dev-mobile/smtp-server/internal/environment"
	"github.com/a-dev-mobile/smtp-server/internal/handlers/send"
	"github.com/a-dev-mobile/smtp-server/internal/logging"
	"github.com/a-dev-mobile/smtp-server/internal/models"
	pb "github.com/a-dev-mobile/smtp-server/proto"
	"google.golang.org/grpc"

	"golang.org/x/exp/slog"
)



func main() {
	appEnv := getAppEnvOrFail()



	cfg, lg := getConfigAndLogOrFail()

	lg.Info("Environment used", ".env", appEnv)
	lg.Debug("init SMTPServer", "config_json", cfg)


	var opts []grpc.ServerOption

	if cfg.GRPCServer.MaxConcurrentStreams > 0 {
		opts = append(opts, grpc.MaxConcurrentStreams(uint32(cfg.GRPCServer.MaxConcurrentStreams)))
	}
// Initialize and start the gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServer.Port))
	if err != nil {
		lg.Error("Failed to listen:", logging.Err(err))
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(opts...)

	// Register services
	sendServiceServer := send.NewServiceServer(cfg, lg)
	pb.RegisterEmailSenderApiServer(grpcServer, sendServiceServer)

	lg.Info("gRPC server starting", "port", cfg.GRPCServer.Port)
	if err := grpcServer.Serve(lis); err != nil {
		lg.Error("Failed to serve:", logging.Err(err))
		os.Exit(1)
	}
}

func getAppEnvOrFail() environment.Environment {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		log.Fatalf("APP_ENV is not set")
	}
	return environment.Environment(appEnv)
}


func getConfigAndLogOrFail() (*models.Config, *slog.Logger) {
	var cfg *models.Config
    var err error

	cfg, err = config.GetConfig[models.Config]("../config","config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	lg := logging.SetupLogger(cfg.Logging.Level, cfg.Logging.FileOutput.FilePath)

	return cfg, lg
}
