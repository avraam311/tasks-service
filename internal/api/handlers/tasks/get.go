package tasks

import (
	"log/slog"
	"net/http"

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
