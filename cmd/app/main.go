package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"time"
	"webhook-server-test/internal"
	"webhook-server-test/internal/service"
)

func main() {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel))
	if err != nil {
		log.Fatalf("zap logger initialise failed: %v", err)
	}
	defer logger.Sync()

	client := &http.Client{Timeout: 5 * time.Second}
	svc := service.NewService(logger, client)
	server := internal.NewServer(":8888", logger, svc)

	server.Start()
}
