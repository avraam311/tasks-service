package tasks

import (
	"context"

	"github.com/avraam311/tasks-service/internal/models"
)

func (r *Repo) LoadTask(ctx context.Context, taskID uint) (*models.TaskDomain, error) {
	r.mu.RLock()
	task, ok := r.storage[taskID]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrTaskNotFound
	}
	taskDomain := models.TaskDomain{
		ID:          taskID,
		Header:      task.Header,
		Description: task.Description,
		Finished:    task.Finished,
	}

	return &taskDomain, nil
}
