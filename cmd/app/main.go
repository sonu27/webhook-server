package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"webhook-server-test/internal"
	"webhook-server-test/internal/service"
)

func main() {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel))
	if err != nil {
		log.Fatalf("zap logger initialise failed: %v", err)
	}
	defer logger.Sync()

	svc := service.NewService(logger)
	server := internal.NewServer(":8888", logger, svc)

	server.StartWithGracefulShutdown()
}
