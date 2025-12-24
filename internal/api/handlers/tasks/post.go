package tasks

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/avraam311/tasks-service/internal/api/responses"
	"github.com/avraam311/tasks-service/internal/models"
)

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Error("failed to decode JSON", slog.String("method", r.Method))
		err := responses.ResponseError(w, responses.ErrMethodNotAllowed, "only POST allowed", http.StatusMethodNotAllowed)
		if err != nil {
			slog.Error("failed to send json response", slog.Any("err", err))
		}
		return
	}

	var task models.TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		slog.Error("failed to decode JSON", slog.Any("error", err), slog.Any("task", r.Body))
		err := responses.ResponseError(w, responses.ErrInvalidJSON, fmt.Sprintf("invalid request body: %s", err.Error()),
			http.StatusBadRequest)
		if err != nil {
			slog.Error("failed to send json response", slog.Any("err", err))
		}
		return
	}

	taskID, err := h.service.CreateTask(r.Context(), &task)
	if err != nil {
		slog.Error("failed to create task", slog.Any("task", task), slog.Any("error", err))
		err := responses.ResponseError(w, responses.ErrInternalServer, "internal server error", http.StatusInternalServerError)
		if err != nil {
			slog.Error("failed to send json response", slog.Any("err", err))
		}
		return
	}

	responses.ResponseCreated(w, taskID)
}
