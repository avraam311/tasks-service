package tasks

import (
	"context"

	"github.com/avraam311/tasks-service/internal/models"
)

func (r *Repo) StoreTask(ctx context.Context, task *models.TaskDTO) (uint, error) {
	taskID := r.taskID
	r.mu.Lock()
	r.storage[r.taskID] = task
	r.taskID++
	r.mu.Unlock()

	return taskID, nil
}
