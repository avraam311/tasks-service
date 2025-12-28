package tasks

import (
	"context"

	"github.com/avraam311/tasks-service/internal/models"
)

type Service interface {
	CreateTask(ctx context.Context, task *models.TaskDTO) (uint, error)
	GetAllTasks(ctx context.Context) ([]*models.TaskDomain, error)
	GetTask(ctx context.Context, taskID uint) (*models.TaskDomain, error)
	UpdateTask(ctx context.Context, taskID uint, task *models.TaskDTO) error
	DeleteTask(ctx context.Context, taskID uint) error
}

type Handler struct {
	service Service
}

func New(service Service) Handler {
	return Handler{
		service: service,
	}
}
