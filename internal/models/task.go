package models

import "time"

// TaskStatus represents the current state of a task.
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

// TaskType represents the kind of task to be executed.
type TaskType string

const (
	TypePrintMessage TaskType = "print_message"
	TypeProcessImage TaskType = "process_image"
	TypeSendEmail    TaskType = "send_email"
)

// Task represents a unit of work to be executed.
type Task struct {
	ID        string     `json:"id"`
	Type      TaskType   `json:"type"`
	Payload   []byte     `json:"payload"`
	Status    TaskStatus `json:"status"`
	Result    string     `json:"result,omitempty"`
	Error     string     `json:"error,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// PrintMessagePayload represents the data needed for a print message task.
type PrintMessagePayload struct {
	Message string `json:"message"`
}

// ProcessImagePayload represents the data needed for an image processing task.
type ProcessImagePayload struct {
	ImageURL string `json:"image_url"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// SendEmailPayload represents the data needed for an email task.
type SendEmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
