package server

import (
	"net/http"

	"github.com/avraam311/tasks-service/internal/api/middlewares"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	router := middlewares.RecoveryMiddleware(mux)
	router = middlewares.LoggingMiddleware(router)

	return router
}

func NewServer(addr string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
