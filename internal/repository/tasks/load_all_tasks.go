package tasks

import (
	"context"

	"github.com/avraam311/tasks-service/internal/models"
)

func (r *Repo) LoadAllTasks(ctx context.Context) ([]*models.TaskDomain, error) {
	tasks := []*models.TaskDomain{}
	r.mu.RLock()
	for taskID, task := range r.storage {
		taskDomain := models.TaskDomain{
			ID:          taskID,
			Header:      task.Header,
			Description: task.Description,
			Finished:    task.Finished,
		}
		tasks = append(tasks, &taskDomain)
	}
	r.mu.RUnlock()

	return tasks, nil
}
