package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlerTasks "github.com/avraam311/tasks-service/internal/api/handlers/tasks"
	"github.com/avraam311/tasks-service/internal/api/server"
	"github.com/avraam311/tasks-service/internal/infra/config"
	"github.com/avraam311/tasks-service/internal/infra/logger"
	repoTasks "github.com/avraam311/tasks-service/internal/repository/tasks"
	serviceTasks "github.com/avraam311/tasks-service/internal/service/tasks"
)

const (
	loggerLevel = "info"
	loggerJSON  = true
	configPath  = "./config/local.json"
)

func main() {
	logger.Init(loggerLevel, loggerJSON)
	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to init config ", "error", err)
		os.Exit(1)
	}
	err = cfg.LoadJSON(configPath)
	if err != nil {
		slog.Error("failed to init config", "error", err)
		os.Exit(1)
	}

	repo := repoTasks.New()
	service := serviceTasks.New(repo)
	handler := handlerTasks.New(service)

	router := server.NewRouter(handler)
	srv := server.NewServer(cfg.Server.Port, router)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("failed to run server", "error", err)
			os.Exit(1)
		}
	}()
	slog.Info("server is running")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()
	slog.Info("shutdown signal recieved")

	shutdownCtx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	slog.Info("shutting down")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Warn("failed to shutdown server", "error", err)
	}
	if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
		slog.Info("timeout exceeded, forcing shutdown")
	}
}
