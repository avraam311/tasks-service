package tasks

import (
	"context"
	"fmt"

	"github.com/avraam311/tasks-service/internal/models"
)

func (s *Service) UpdateTask(ctx context.Context, taskID uint, task *models.TaskDTO) error {
	err := s.repo.SwapTask(ctx, taskID, task)
	if err != nil {
		return fmt.Errorf("service/update_task.go - %w", err)
	}

	return nil
}
