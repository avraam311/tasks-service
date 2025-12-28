package tasks

import (
	"context"
	"fmt"

	"github.com/avraam311/tasks-service/internal/models"
)

func (s *Service) GetAllTasks(ctx context.Context) ([]*models.TaskDomain, error) {
	tasks, err := s.repo.LoadAllTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("service/get_all_tasks.go - %w", err)
	}

	return tasks, nil
}
