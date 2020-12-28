package internal

import (
	"net/http"
	"webhook-server-test/internal/service"
)

func NewServer(addr string, svc *service.Service) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/webhooks", http.HandlerFunc(svc.CreateWebhooksHandler))
	mux.Handle("/fire-webhooks", http.HandlerFunc(svc.FireWebhooksHandler))

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
