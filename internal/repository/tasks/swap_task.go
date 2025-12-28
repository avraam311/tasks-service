package tasks

import (
	"context"

	"github.com/avraam311/tasks-service/internal/models"
)

func (r *Repo) SwapTask(ctx context.Context, taskID uint, task *models.TaskDTO) error {
	r.mu.Lock()
	r.storage[taskID] = task
	r.mu.Unlock()

	return nil
}
