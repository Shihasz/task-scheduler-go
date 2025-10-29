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
	log.Printf("Executing print message task with payload: %s", string(task.Payload))

	// The payload should be JSON bytes, we need to unmarshal it.
	var payload models.PrintMessagePayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		log.Printf("Failed to unmarshal print message payload: %v, payload: %s", err, string(task.Payload))
		return "", err
	}

	// Simulate task execution.
	log.Printf("PRINT_MESSAGE: %s", payload.Message)

	// Simulate some processing time.
	// time.Sleep(100 * time.Millisecond)

	return "Message printed successfully: " + payload.Message, nil
}

// CanHandle returns true if this executor can handle the task type.
func (e *PrintMessageExecutor) CanHandle(taskType models.TaskType) bool {
	return taskType == models.TypePrintMessage
}
