package storage

import "task-scheduler-go/internal/models"

// Storage interface defines the contract for task storage.
type Storage interface {
	CreateTask(taskType models.TaskType, payload []byte) (*models.Task, error)
	GetTask(id string) (*models.Task, error)
	UpdateTask(id string, status models.TaskStatus, result, errorMsg string) error
	ListTasks(status models.TaskStatus) ([]*models.Task, error)
}
