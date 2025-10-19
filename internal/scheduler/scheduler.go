package scheduler

import (
	"task-scheduler-go/internal/models"
	"task-scheduler-go/internal/storage"
)

// Scheduler manages task queue and distribution.
type Scheduler struct {
	storage storage.Storage
}

// NewScheduler creates a new scheduler instance.
func NewScheduler(storage storage.Storage) *Scheduler {
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
