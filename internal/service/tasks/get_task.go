package tasks

import (
	"context"
	"fmt"

	"github.com/avraam311/tasks-service/internal/models"
)

func (s *Service) GetTask(ctx context.Context, taskID uint) (*models.TaskDomain, error) {
	task, err := s.repo.LoadTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("service/get_task.go - %w", err)
	}

	return task, nil
}
