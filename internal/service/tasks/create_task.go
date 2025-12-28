package tasks

import (
	"context"
	"fmt"

	"github.com/avraam311/tasks-service/internal/models"
)

func (s *Service) CreateTask(ctx context.Context, task *models.TaskDTO) (uint, error) {
	taskID, err := s.repo.StoreTask(ctx, task)
	if err != nil {
		return 0, fmt.Errorf("service/create_task.go - %w", err)
	}

	return taskID, nil
}
