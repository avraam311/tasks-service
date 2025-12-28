package tasks

import (
	"context"
	"fmt"
)

func (s *Service) DeleteTask(ctx context.Context, taskID uint) error {
	err := s.repo.DeleteTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("service/delete_task.go - %w", err)
	}

	return nil
}
