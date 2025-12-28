package tasks

import (
	"context"

	"github.com/avraam311/tasks-service/internal/models"
)

type Repo interface {
	StoreTask(ctx context.Context, task *models.TaskDTO) (uint, error)
	LoadAllTasks(ctx context.Context) ([]*models.TaskDomain, error)
}

type Service struct {
	repo Repo
}

func New(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}
