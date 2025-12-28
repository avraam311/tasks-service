package tasks

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avraam311/tasks-service/internal/models"
)

func TestRepo_New(t *testing.T) {
	repo := New()

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.storage)
	assert.Equal(t, uint(0), repo.taskID)
}

func TestRepo_StoreTask(t *testing.T) {
	ctx := context.Background()

	t.Run("First Task", func(t *testing.T) {
		repo := New()
		task := &models.TaskDTO{
			Header:      "Test Task 1",
			Description: "Test Description 1",
			Finished:    false,
		}

		taskID, err := repo.StoreTask(ctx, task)

		assert.NoError(t, err)
		assert.Equal(t, uint(0), taskID)

		storedTask, ok := repo.storage[taskID]
		assert.True(t, ok)
		assert.Equal(t, task, storedTask)
	})

	t.Run("Second Task", func(t *testing.T) {
		repo := New()
		task := &models.TaskDTO{
			Header:      "Test Task 2",
			Description: "Test Description 2",
			Finished:    true,
		}

		taskID, err := repo.StoreTask(ctx, task)

		assert.NoError(t, err)
		assert.Equal(t, uint(0), taskID)

		storedTask, ok := repo.storage[taskID]
		assert.True(t, ok)
		assert.Equal(t, task, storedTask)
	})

	t.Run("Auto-increment", func(t *testing.T) {
		repo := New()

		task1 := &models.TaskDTO{Header: "Task 1"}
		task2 := &models.TaskDTO{Header: "Task 2"}
		task3 := &models.TaskDTO{Header: "Task 3"}

		id1, _ := repo.StoreTask(ctx, task1)
		id2, _ := repo.StoreTask(ctx, task2)
		id3, _ := repo.StoreTask(ctx, task3)

		assert.Equal(t, uint(0), id1)
		assert.Equal(t, uint(1), id2)
		assert.Equal(t, uint(2), id3)
		assert.Equal(t, uint(3), repo.taskID)
	})
}

func TestRepo_LoadTask(t *testing.T) {
	ctx := context.Background()

	t.Run("Existing Task", func(t *testing.T) {
		repo := New()

		task := &models.TaskDTO{
			Header:      "Test Task",
			Description: "Test Description",
			Finished:    false,
		}
		taskID, _ := repo.StoreTask(ctx, task)

		loadedTask, err := repo.LoadTask(ctx, taskID)

		assert.NoError(t, err)
		assert.NotNil(t, loadedTask)
		assert.Equal(t, taskID, loadedTask.ID)
		assert.Equal(t, task.Header, loadedTask.Header)
		assert.Equal(t, task.Description, loadedTask.Description)
		assert.Equal(t, task.Finished, loadedTask.Finished)
	})

	t.Run("Non-existent Task", func(t *testing.T) {
		repo := New()
		nonExistentID := uint(999)

		loadedTask, err := repo.LoadTask(ctx, nonExistentID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrTaskNotFound))
		assert.Nil(t, loadedTask)
	})

	t.Run("Zero ID in Empty Repo", func(t *testing.T) {
		repo := New()

		loadedTask, err := repo.LoadTask(ctx, uint(0))

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrTaskNotFound))
		assert.Nil(t, loadedTask)
	})
}

func TestRepo_LoadAllTasks(t *testing.T) {
	ctx := context.Background()

	t.Run("Empty Repository", func(t *testing.T) {
		repo := New()

		tasks, err := repo.LoadAllTasks(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, tasks)
		assert.Len(t, tasks, 0)
	})

	t.Run("Single Task", func(t *testing.T) {
		repo := New()

		task := &models.TaskDTO{
			Header:      "Single Task",
			Description: "Single Description",
			Finished:    true,
		}
		taskID, _ := repo.StoreTask(ctx, task)

		tasks, err := repo.LoadAllTasks(ctx)

		assert.NoError(t, err)
		assert.Len(t, tasks, 1)

		loadedTask := tasks[0]
		assert.Equal(t, taskID, loadedTask.ID)
		assert.Equal(t, task.Header, loadedTask.Header)
		assert.Equal(t, task.Description, loadedTask.Description)
		assert.Equal(t, task.Finished, loadedTask.Finished)
	})

	t.Run("Multiple Tasks", func(t *testing.T) {
		repo := New()

		task1 := &models.TaskDTO{Header: "Task 1", Description: "Desc 1", Finished: false}
		task2 := &models.TaskDTO{Header: "Task 2", Description: "Desc 2", Finished: true}
		task3 := &models.TaskDTO{Header: "Task 3", Description: "Desc 3", Finished: false}

		id1, _ := repo.StoreTask(ctx, task1)
		id2, _ := repo.StoreTask(ctx, task2)
		id3, _ := repo.StoreTask(ctx, task3)

		tasks, err := repo.LoadAllTasks(ctx)

		assert.NoError(t, err)
		assert.Len(t, tasks, 3)

		taskMap := make(map[uint]*models.TaskDomain)
		for _, task := range tasks {
			taskMap[task.ID] = task
		}

		assert.Contains(t, taskMap, id1)
		assert.Contains(t, taskMap, id2)
		assert.Contains(t, taskMap, id3)

		assert.Equal(t, task1.Header, taskMap[id1].Header)
		assert.Equal(t, task2.Header, taskMap[id2].Header)
		assert.Equal(t, task3.Header, taskMap[id3].Header)
	})
}

func TestRepo_SwapTask(t *testing.T) {
	ctx := context.Background()

	t.Run("Update Existing Task", func(t *testing.T) {
		repo := New()

		originalTask := &models.TaskDTO{
			Header:      "Original Task",
			Description: "Original Description",
			Finished:    false,
		}
		taskID, _ := repo.StoreTask(ctx, originalTask)

		updatedTask := &models.TaskDTO{
			Header:      "Updated Task",
			Description: "Updated Description",
			Finished:    true,
		}

		err := repo.SwapTask(ctx, taskID, updatedTask)
		assert.NoError(t, err)

		loadedTask, err := repo.LoadTask(ctx, taskID)
		assert.NoError(t, err)
		assert.Equal(t, updatedTask.Header, loadedTask.Header)
		assert.Equal(t, updatedTask.Description, loadedTask.Description)
		assert.Equal(t, updatedTask.Finished, loadedTask.Finished)
	})

	t.Run("Create New Task", func(t *testing.T) {
		repo := New()

		newTask := &models.TaskDTO{
			Header:      "New Task",
			Description: "New Description",
			Finished:    false,
		}
		newID := uint(42)

		err := repo.SwapTask(ctx, newID, newTask)
		assert.NoError(t, err)

		loadedTask, err := repo.LoadTask(ctx, newID)
		assert.NoError(t, err)
		assert.Equal(t, newTask.Header, loadedTask.Header)
		assert.Equal(t, newTask.Description, loadedTask.Description)
		assert.Equal(t, newTask.Finished, loadedTask.Finished)
	})
}

func TestRepo_DeleteTask(t *testing.T) {
	ctx := context.Background()

	t.Run("Delete Existing Task", func(t *testing.T) {
		repo := New()

		task := &models.TaskDTO{
			Header:      "Task to Delete",
			Description: "Description",
			Finished:    false,
		}
		taskID, _ := repo.StoreTask(ctx, task)

		err := repo.DeleteTask(ctx, taskID)
		assert.NoError(t, err)

		loadedTask, err := repo.LoadTask(ctx, taskID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrTaskNotFound))
		assert.Nil(t, loadedTask)
	})

	t.Run("Delete Non-existent Task", func(t *testing.T) {
		repo := New()

		err := repo.DeleteTask(ctx, uint(999))
		assert.NoError(t, err)
	})
}

func TestRepo_Integration(t *testing.T) {
	ctx := context.Background()
	repo := New()

	createTask := &models.TaskDTO{
		Header:      "Integration Task",
		Description: "Integration Description",
		Finished:    false,
	}

	taskID, err := repo.StoreTask(ctx, createTask)
	require.NoError(t, err)
	assert.Equal(t, uint(0), taskID)

	loadedTask, err := repo.LoadTask(ctx, taskID)
	require.NoError(t, err)
	assert.Equal(t, createTask.Header, loadedTask.Header)

	updatedTask := &models.TaskDTO{
		Header:      "Updated Integration Task",
		Description: "Updated Integration Description",
		Finished:    true,
	}
	err = repo.SwapTask(ctx, taskID, updatedTask)
	require.NoError(t, err)

	loadedUpdatedTask, err := repo.LoadTask(ctx, taskID)
	require.NoError(t, err)
	assert.Equal(t, updatedTask.Header, loadedUpdatedTask.Header)
	assert.Equal(t, updatedTask.Finished, loadedUpdatedTask.Finished)

	err = repo.DeleteTask(ctx, taskID)
	require.NoError(t, err)

	_, err = repo.LoadTask(ctx, taskID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrTaskNotFound))
}

func TestRepo_ThreadSafety(t *testing.T) {
	ctx := context.Background()
	repo := New()

	const numGoroutines = 10
	const numOperations = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				task := &models.TaskDTO{
					Header:      fmt.Sprintf("Goroutine %d Task %d", goroutineID, j),
					Description: fmt.Sprintf("Description from goroutine %d", goroutineID),
					Finished:    j%2 == 0,
				}

				taskID, err := repo.StoreTask(ctx, task)
				assert.NoError(t, err)

				loadedTask, err := repo.LoadTask(ctx, taskID)
				assert.NoError(t, err)
				assert.NotNil(t, loadedTask)
				assert.Equal(t, taskID, loadedTask.ID)
				assert.NotEmpty(t, loadedTask.Header)
				assert.NotEmpty(t, loadedTask.Description)
				assert.IsType(t, bool(false), loadedTask.Finished)
			}
		}(i)
	}

	wg.Wait()

	allTasks, err := repo.LoadAllTasks(ctx)
	require.NoError(t, err)

	expectedCount := numGoroutines * numOperations
	assert.Len(t, allTasks, expectedCount)

	assert.Equal(t, uint(expectedCount), repo.taskID)

	taskMap := make(map[uint]*models.TaskDomain)
	for _, task := range allTasks {
		assert.NotContains(t, taskMap, task.ID)
		taskMap[task.ID] = task

		assert.NotEmpty(t, task.Header)
		assert.NotEmpty(t, task.Description)
		assert.Greater(t, len(task.Header), 0)
		assert.Greater(t, len(task.Description), 0)
		assert.IsType(t, bool(false), task.Finished)
	}

	for i := uint(0); i < uint(expectedCount); i++ {
		assert.Contains(t, taskMap, i, "Task with ID %d should exist", i)
	}
}
