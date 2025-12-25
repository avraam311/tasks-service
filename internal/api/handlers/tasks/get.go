package tasks

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/avraam311/tasks-service/internal/api/responses"
)

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		slog.Error("not allowed method", slog.String("method", r.Method))
		err := responses.ResponseError(w, responses.ErrMethodNotAllowed, "only GET allowed", http.StatusMethodNotAllowed)
		if err != nil {
			slog.Error("failed to send json response", slog.Any("err", err))
		}
		return
	}

	tasks, err := h.service.GetAllTasks(r.Context())
	if err != nil {
		slog.Error("failed to get all tasks", slog.Any("error", err))
		err := responses.ResponseError(w, responses.ErrInternalServer, "internal server error", http.StatusInternalServerError)
		if err != nil {
			slog.Error("failed to send json response", slog.Any("err", err))
		}
		return
	}

	responses.ResponseOK(w, tasks)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		slog.Error("not allowed method", slog.String("method", r.Method))
		err := responses.ResponseError(w, responses.ErrMethodNotAllowed, "only GET allowed", http.StatusMethodNotAllowed)
		if err != nil {
			slog.Error("failed to send json response", slog.Any("err", err))
		}
		return
	}

	taskIDStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	taskIDInt, err := strconv.Atoi(taskIDStr)
	if err != nil {
		slog.Error("failed to convert task id into int", slog.String("tasks id str", taskIDStr))
		err := responses.ResponseError(w, responses.ErrInvalidID, "invalid task id", http.StatusBadRequest)
		if err != nil {
			slog.Error("failed to send json response", slog.Any("err", err))
		}
		return
	}
	taskID := uint(taskIDInt)

	task, err := h.service.GetTask(r.Context(), taskID)
	if err != nil {
		slog.Error("failed to get task", slog.Any("error", err))
		err := responses.ResponseError(w, responses.ErrInternalServer, "internal server error", http.StatusInternalServerError)
		if err != nil {
			slog.Error("failed to send json response", slog.Any("err", err))
		}
		return
	}

	responses.ResponseOK(w, task)
}
