package tasks

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/avraam311/tasks-service/internal/api/responses"
	"github.com/avraam311/tasks-service/internal/repository/tasks"
)

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		slog.Error("not allowed method", slog.String("method", r.Method))
		err := responses.ResponseError(w, responses.ErrMethodNotAllowed, "only DELETE allowed", http.StatusMethodNotAllowed)
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

	err = h.service.DeleteTask(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, tasks.ErrTaskNotFound) {
			slog.Error("task not found", slog.Any("task_id", taskID))
			err := responses.ResponseError(w, responses.ErrTaskNotFound, "task not found", http.StatusBadRequest)
			if err != nil {
				slog.Error("failed to send json response", slog.Any("err", err))
			}
			return
		}

		slog.Error("failed to delete task", slog.Any("task id", taskID), slog.Any("error", err))
		err := responses.ResponseError(w, responses.ErrInternalServer, "internal server error", http.StatusInternalServerError)
		if err != nil {
			slog.Error("failed to send json response", slog.Any("err", err))
		}
		return
	}

	err = responses.ResponseOK(w, responses.SuccessTaskDeleted)
	if err != nil {
		slog.Error("failed to send json response", slog.Any("err", err))
	}
}
