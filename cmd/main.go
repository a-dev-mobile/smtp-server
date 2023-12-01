package main

import (
	"fmt"

	"github.com/a-dev-mobile/smtp-server/internal/handlers/send"
	"log"
	"net"
	"os"

	pb "github.com/a-dev-mobile/smtp-server/proto"

	"github.com/a-dev-mobile/smtp-server/internal/config"


	"golang.org/x/exp/slog"
	"google.golang.org/grpc"

	"github.com/a-dev-mobile/smtp-server/lib/logger/sl"
)



func main() {
	 cfg, lg := initializeApp()

	lg.Info("init SMTPServer", "config_json", cfg)


	var opts []grpc.ServerOption

	if cfg.GRPCServer.MaxConcurrentStreams > 0 {
		opts = append(opts, grpc.MaxConcurrentStreams(uint32(cfg.GRPCServer.MaxConcurrentStreams)))
	}
	// Инициализация и запуск gRPC сервера
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServer.Port))
	if err != nil {
		lg.Error("Failed to listen:", sl.Err(err))
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(opts...)

	// Регистрация сервисов
	sendServiceServer := send.NewServiceServer(cfg, lg)
	pb.RegisterEmailSenderApiServer(grpcServer, sendServiceServer)

	lg.Info("gRPC server starting", "port", cfg.GRPCServer.Port)
	if err := grpcServer.Serve(lis); err != nil {
		lg.Error("Failed to serve:", sl.Err(err))
		os.Exit(1)
	}
}

// initializeApp sets up the application environment, configuration, and logger.
// It determines the application's running environment, loads the appropriate configuration,
// and initializes the logging system.
func initializeApp() ( *config.Config, *slog.Logger) {

	cfg, lg := getConfigOrFail()
	return  cfg, lg
}



func getConfigOrFail() (*config.Config, *slog.Logger) {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	lg := sl.SetupLogger(cfg.Logging.Level, cfg.Logging.FileOutput.FilePath)

	return cfg, lg
}
