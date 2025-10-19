package storage

import (
	"sync"
	"task-scheduler-go/internal/models"
	"time"

	"github.com/google/uuid"
)

// MemoryStorage implements an in-memory task store.
type MemoryStorage struct {
	tasks map[string]*models.Task
	mu    sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks: make(map[string]*models.Task),
	}
}

// CreateTask stores a new task.
func (s *MemoryStorage) CreateTask(taskType models.TaskType, payload []byte) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := &models.Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Payload:   payload,
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.tasks[task.ID] = task
	return task, nil
}

// GetTask retrieves a task by its ID.
func (s *MemoryStorage) GetTask(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// UpdateTask updates a task's status and result.
func (s *MemoryStorage) UpdateTask(id string, status models.TaskStatus, result, errorMsg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return ErrTaskNotFound
	}

	task.Status = status
	task.Result = result
	task.Error = errorMsg
	task.UpdatedAt = time.Now()

	return nil
}

// ListTasks returns all tasks, optionally filtered by status.
func (s *MemoryStorage) ListTasks(status models.TaskStatus) ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tasks []*models.Task
	for _, task := range s.tasks {
		if status == "" || task.Status == status {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}
