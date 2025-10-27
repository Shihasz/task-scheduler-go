package worker

import (
	"task-scheduler-go/internal/models"
)

// ExecutorRegistry manages all task executors.
type ExecutorRegistry struct {
	executors []TaskExecutor
}

// NewExecutorRegistry creates a new executor registry.
func NewExecutorRegistry() *ExecutorRegistry {
	return &ExecutorRegistry{
		executors: make([]TaskExecutor, 0),
	}
}

// Register adds a new executor to the registry.
func (r *ExecutorRegistry) Register(executor TaskExecutor) {
	r.executors = append(r.executors, executor)
}

// GetExecutorForTask finds an executor that can handle the given task type.
func (r *ExecutorRegistry) GetExecutorForTask(taskType models.TaskType) TaskExecutor {
	for _, executor := range r.executors {
		if executor.CanHandle(taskType) {
			return executor
		}
	}
	return nil
}

// GetExecutors returns all registered executors.
func (r *ExecutorRegistry) GetExecutors() []TaskExecutor {
	return r.executors
}
