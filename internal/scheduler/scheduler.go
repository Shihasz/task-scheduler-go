package scheduler

import (
	"task-scheduler-go/internal/models"
	"task-scheduler-go/internal/storage"
)

// Scheduler manages task queue and distribution.
type Scheduler struct {
	storage storage.Storage
}

// Storage interface defines the contract for task storage.
type Storage interface {
	CreateTask(taskType models.TaskType, payload []byte) (*models.Task, error)
	GetTask(id string) (*models.Task, error)
	UpdateTask(id string, status models.TaskStatus, result, errorMsg string) error
	ListTasks(status models.TaskStatus) ([]*models.Task, error)
}

// NewScheduler creates a new scheduler instance.
func NewScheduler(storage Storage) *Scheduler {
	return &Scheduler{
		storage: storage,
	}
}

// SubmitTask creates a new task and adds it to the queue.
func (s *Scheduler) SubmitTask(taskType models.TaskType, payload []byte) (*models.Task, error) {
	return s.storage.CreateTask(taskType, payload)
}

// GetTaskStatus retrieves the current status of a task.
func (s *Scheduler) GetTaskStatus(taskID string) (*models.Task, error) {
	return s.storage.GetTask(taskID)
}

// ListPendingTasks returns all pending tasks.
func (s *Scheduler) ListPendingTasks() ([]*models.Task, error) {
	return s.storage.ListTasks(models.StatusPending)
}
