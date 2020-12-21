package internal

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	*http.Server
	l *zap.Logger
}

func NewServer(addr string, logger *zap.Logger) *Server {
	mux := http.NewServeMux()
	mux.Handle("/webhooks", Recovery(http.HandlerFunc(createWebhooksHandler)))
	mux.Handle("/fire-webhooks", Recovery(http.HandlerFunc(testWebhooksHandler)))

	server := http.Server{
		Handler: mux,
		Addr:    addr,
	}

	return &Server{
		l:      logger,
		Server: &server,
	}
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println(err)

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
