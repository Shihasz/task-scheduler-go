package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-scheduler-go/internal/models"
	"task-scheduler-go/internal/scheduler"
	"task-scheduler-go/internal/storage"
	"task-scheduler-go/internal/worker"
	"task-scheduler-go/internal/worker/executors"
)

func main() {
	// Use the same storage instance for both scheduler and worker
	store := storage.NewMemoryStorage()

	// Start Scheduler
	sched := scheduler.NewScheduler(store)

	// Start Worker
	workerID := "worker-1"
	w := worker.NewWorker(workerID, store)

	// Initialize executors
	registry := worker.NewExecutorRegistry()
	registry.Register(executors.NewPrintMessageExecutor())

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Printf("Shutdown signal received, stopping worker...")
		w.Stop()
		os.Exit(0)
	}()

	// Start the worker in a goroutine
	go func() {
		log.Printf("Starting worker %s...", workerID)
		w.Start(registry.GetExecutors())
	}()

	// HTTP handlers.
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateTask(sched, w, r)
		case http.MethodGet:
			handleListTasks(sched, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetTask(sched, w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Debug endpoint to see all tasks
	http.HandleFunc("/debug/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := store.ListTasks("")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	})

	// Start server
	port := ":8080"
	log.Printf("Task Scheduler Server starting on port %s", port)
	log.Printf("API Endpoints:")
	log.Printf("	POST   /tasks         - Submit a new task")
	log.Printf("	GET    /tasks         - List pending tasks")
	log.Printf("	GET    /tasks/{id}    - Get task status")
	log.Printf("	GET    /debug/tasks   - Debug all tasks")
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleCreateTask(sched *scheduler.Scheduler, w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type    models.TaskType `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log the received payload for debugging
	log.Printf("Received task creation request: type=%s, payload=%s", req.Type, string(req.Payload))

	task, err := sched.SubmitTask(req.Type, req.Payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func handleGetTask(sched *scheduler.Scheduler, w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/tasks/"):]
	if taskID == "" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}

	task, err := sched.GetTaskStatus(taskID)
	if err != nil {
		if err == storage.ErrTaskNotFound {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func handleListTasks(sched *scheduler.Scheduler, w http.ResponseWriter, r *http.Request) {
	tasks, err := sched.ListPendingTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
