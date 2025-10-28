package worker

import (
	"log"
	"task-scheduler-go/internal/models"
	"task-scheduler-go/internal/storage"
	"time"
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

// Start begins the worker's  task polling and execution.
func (w *Worker) Start(executors []TaskExecutor) {
	log.Printf("Worker %s starting with %d executors...", w.ID, len(executors))

	ticker := time.NewTicker(2 * time.Second) // Poll every 2 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.pollAndExecute(executors)
		case <-w.stopChan:
			log.Printf("Worker %s stopping...", w.ID)
			return
		}
	}
}

// Stop signals the worker to stop processing.
func (w *Worker) Stop() {
	close(w.stopChan)
}

// pollAndExecute fetches pending tasks and executes them.
func (w *Worker) pollAndExecute(executors []TaskExecutor) {
	// Get pending tasks from storage
	tasks, err := w.storage.ListTasks(models.StatusPending)
	if err != nil {
		log.Printf("Worker %s failed to fetch tasks: %v", w.ID, err)
		return
	}

	if len(tasks) == 0 {
		log.Printf("Worker %s: no pending tasks found", w.ID)
		return // No tasks to process
	}

	log.Printf("Worker %s found %d pending tasks", w.ID, len(tasks))

	// Process each task
	for _, task := range tasks {
		w.executeTask(task, executors)
	}
}

// executeTask handles the execution of a single task.
func (w *Worker) executeTask(task *models.Task, executors []TaskExecutor) {
	log.Printf("Worker %s processing task %s (type: %s)", w.ID, task.ID, task.Type)
	// Update task status to running
	if err := w.storage.UpdateTask(task.ID, models.StatusRunning, "", ""); err != nil {
		log.Printf("Worker %s failed to update task %s to running: %v", w.ID, task.ID, err)
		return
	}

	log.Printf("Worker %s executing task %s (type: %s)", w.ID, task.ID, task.Type)

	// Find appropriate executor
	var executor TaskExecutor
	for _, ex := range executors {
		if ex.CanHandle(task.Type) {
			executor = ex
			break
		}
	}

	if executor == nil {
		// No executor found for this task type
		errMsg := "No executor found for task type: " + string(task.Type)
		w.storage.UpdateTask(task.ID, models.StatusFailed, "", errMsg)
		log.Printf("Worker %s failed to execute task %s: %s", w.ID, task.ID, errMsg)
		return
	}

	// Execute the task
	result, err := executor.Execute(task)
	if err != nil {
		// Task execution failed
		w.storage.UpdateTask(task.ID, models.StatusFailed, "", err.Error())
		log.Printf("Worker %s failed to execute task %s: %v", w.ID, task.ID, err)
		return
	}

	// Task execution succeeded
	w.storage.UpdateTask(task.ID, models.StatusCompleted, result, "")
	log.Printf("Worker %s completed task %s: %s", w.ID, task.ID, result)
}
