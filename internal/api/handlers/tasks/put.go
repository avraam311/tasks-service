package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/avraam311/tasks-service/internal/api/responses"
	"github.com/avraam311/tasks-service/internal/models"
	"github.com/avraam311/tasks-service/internal/repository/tasks"
)

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		slog.Error("not allowed method", slog.String("method", r.Method))
		err := responses.ResponseError(w, responses.ErrMethodNotAllowed, "only PUT allowed", http.StatusMethodNotAllowed)
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

	err = h.service.UpdateTask(r.Context(), taskID, &task)
	if err != nil {
		if errors.Is(err, tasks.ErrTaskNotFound) {
			slog.Error("task not found", slog.Any("task_id", taskID), slog.Any("task", task))
			err := responses.ResponseError(w, responses.ErrTaskNotFound, "task not found", http.StatusBadRequest)
			if err != nil {
				slog.Error("failed to send json response", slog.Any("err", err))
			}
			return
		}

		slog.Error("failed to update task", slog.Any("task id", taskID), slog.Any("task", task), slog.Any("error", err))
		err := responses.ResponseError(w, responses.ErrInternalServer, "internal server error", http.StatusInternalServerError)
		if err != nil {
			slog.Error("failed to send json response", slog.Any("err", err))
		}
		return
	}

	err = responses.ResponseCreated(w, responses.SuccessTaskUpdated)
	if err != nil {
		slog.Error("failed to send json response", slog.Any("err", err))
	}
}
