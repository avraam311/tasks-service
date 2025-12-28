package tasks

import (
	"errors"
	"sync"

	"github.com/avraam311/tasks-service/internal/models"
)

var (
	ErrTaskNotFound = errors.New("task not found")
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
