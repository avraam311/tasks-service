package tasks

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avraam311/tasks-service/internal/mocks"
	"github.com/avraam311/tasks-service/internal/models"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepo(ctrl)
	service := New(mockRepo)

	testTask := &models.TaskDTO{
		Header:      "Test Task",
		Description: "Test Description",
		Finished:    false,
	}

	tests := []struct {
		name          string
		task          *models.TaskDTO
		repoReturnID  uint
		repoReturnErr error
		expectedID    uint
		expectedErr   string
	}{
		{
			name:          "Success",
			task:          testTask,
			repoReturnID:  123,
			repoReturnErr: nil,
			expectedID:    123,
			expectedErr:   "",
		},
		{
			name:          "RepositoryError",
			task:          testTask,
			repoReturnID:  0,
			repoReturnErr: assert.AnError,
			expectedID:    0,
			expectedErr:   "service/create_task.go -",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo.EXPECT().StoreTask(ctx, tt.task).Return(tt.repoReturnID, tt.repoReturnErr)

			taskID, err := service.CreateTask(ctx, tt.task)

			if tt.expectedErr == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, taskID)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepo(ctrl)
	service := New(mockRepo)

	testTask := &models.TaskDomain{
		ID:          123,
		Header:      "Test Task",
		Description: "Test Description",
		Finished:    false,
	}

	tests := []struct {
		name          string
		taskID        uint
		repoReturn    *models.TaskDomain
		repoReturnErr error
		expectedTask  *models.TaskDomain
		expectedErr   string
	}{
		{
			name:          "Success",
			taskID:        123,
			repoReturn:    testTask,
			repoReturnErr: nil,
			expectedTask:  testTask,
			expectedErr:   "",
		},
		{
			name:          "RepositoryError",
			taskID:        123,
			repoReturn:    nil,
			repoReturnErr: assert.AnError,
			expectedTask:  nil,
			expectedErr:   "service/get_task.go -",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo.EXPECT().LoadTask(ctx, tt.taskID).Return(tt.repoReturn, tt.repoReturnErr)

			task, err := service.GetTask(ctx, tt.taskID)

			if tt.expectedErr == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTask, task)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepo(ctrl)
	service := New(mockRepo)

	testTasks := []*models.TaskDomain{
		{
			ID:          1,
			Header:      "Task 1",
			Description: "Description 1",
			Finished:    false,
		},
		{
			ID:          2,
			Header:      "Task 2",
			Description: "Description 2",
			Finished:    true,
		},
	}

	tests := []struct {
		name          string
		repoReturn    []*models.TaskDomain
		repoReturnErr error
		expectedTasks []*models.TaskDomain
		expectedErr   string
	}{
		{
			name:          "EmptyList",
			repoReturn:    []*models.TaskDomain{},
			repoReturnErr: nil,
			expectedTasks: []*models.TaskDomain{},
			expectedErr:   "",
		},
		{
			name:          "MultipleTasks",
			repoReturn:    testTasks,
			repoReturnErr: nil,
			expectedTasks: testTasks,
			expectedErr:   "",
		},
		{
			name:          "RepositoryError",
			repoReturn:    nil,
			repoReturnErr: assert.AnError,
			expectedTasks: nil,
			expectedErr:   "service/get_all_tasks.go -",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo.EXPECT().LoadAllTasks(ctx).Return(tt.repoReturn, tt.repoReturnErr)

			tasks, err := service.GetAllTasks(ctx)

			if tt.expectedErr == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTasks, tasks)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepo(ctrl)
	service := New(mockRepo)

	testTask := &models.TaskDTO{
		Header:      "Updated Task",
		Description: "Updated Description",
		Finished:    true,
	}

	tests := []struct {
		name          string
		taskID        uint
		task          *models.TaskDTO
		repoReturnErr error
		expectedErr   string
	}{
		{
			name:          "Success",
			taskID:        123,
			task:          testTask,
			repoReturnErr: nil,
			expectedErr:   "",
		},
		{
			name:          "RepositoryError",
			taskID:        123,
			task:          testTask,
			repoReturnErr: assert.AnError,
			expectedErr:   "service/update_task.go -",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo.EXPECT().SwapTask(ctx, tt.taskID, tt.task).Return(tt.repoReturnErr)

			err := service.UpdateTask(ctx, tt.taskID, tt.task)

			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepo(ctrl)
	service := New(mockRepo)

	tests := []struct {
		name          string
		taskID        uint
		repoReturnErr error
		expectedErr   string
	}{
		{
			name:          "Success",
			taskID:        123,
			repoReturnErr: nil,
			expectedErr:   "",
		},
		{
			name:          "RepositoryError",
			taskID:        123,
			repoReturnErr: assert.AnError,
			expectedErr:   "service/delete_task.go -",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo.EXPECT().DeleteTask(ctx, tt.taskID).Return(tt.repoReturnErr)

			err := service.DeleteTask(ctx, tt.taskID)

			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}
