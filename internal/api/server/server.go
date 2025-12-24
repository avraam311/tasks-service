package server

import (
	"net/http"

	"github.com/avraam311/tasks-service/internal/api/handlers/tasks"
	"github.com/avraam311/tasks-service/internal/api/middlewares"
)

func NewRouter(tasksHand tasks.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /todos", tasksHand.CreateTask)

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
