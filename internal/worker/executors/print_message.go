package executors

import (
	"encoding/json"
	"log"
	"task-scheduler-go/internal/models"
)

// PrintMessageExecutor handles print_message tasks.
type PrintMessageExecutor struct{}

// NewPrintMessageExecutor creates a new print message executor.
func NewPrintMessageExecutor() *PrintMessageExecutor {
	return &PrintMessageExecutor{}
}

// Execute processes a print_message task.
func (e *PrintMessageExecutor) Execute(task *models.Task) (string, error) {
	var payload models.PrintMessagePayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return "", err
	}

	// Simulate task execution.
	log.Printf("PRINT_MESSAGE: %s", payload.Message)

	// Simulate some processing time.
	// time.Sleep(100 * time.Millisecond)

	return "Message printed successfully", nil
}

// CanHandle returns true if this executor can handle the task type.
func (e *PrintMessageExecutor) CanHandle(taskType models.TaskType) bool {
	return taskType == models.TypePrintMessage
}
