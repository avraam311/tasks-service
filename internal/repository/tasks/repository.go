package tasks

import (
	"sync"

	"github.com/avraam311/tasks-service/internal/models"
)

type Repo struct {
	storage map[uint]*models.TaskDTO
	taskID  uint
	mu      sync.RWMutex
}

func New() *Repo {
	return &Repo{
		storage: make(map[uint]*models.TaskDTO),
	}
}
