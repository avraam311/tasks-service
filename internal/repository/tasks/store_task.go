package tasks

import (
	"context"

	"github.com/avraam311/tasks-service/internal/models"
)

func (r *Repo) StoreTask(ctx context.Context, task *models.TaskDTO) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage[r.taskID] = task
	r.taskID++
}
