package worker

import (
	"task-scheduler-go/internal/models"
	"task-scheduler-go/internal/storage"
)

// Worker processes tasks from the scheduler.
type Worker struct {
	ID       string
	storage  storage.Storage
	stopChan chan struct{}
}

// NewWorker creates a new Worker instance.
func NewWorker(id string, storage storage.Storage) *Worker {
	return &Worker{
		ID:       id,
		storage:  storage,
		stopChan: make(chan struct{}),
	}
}

// TaskExecutor defines the interface for executing tasks.
type TaskExecutor interface {
	Execute(task *models.Task) (string, error)
	CanHandle(taskType models.TaskType) bool
}
