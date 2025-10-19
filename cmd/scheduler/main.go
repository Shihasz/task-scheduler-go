package main

import (
	"encoding/json"
	"log"
	"net/http"
	"task-scheduler-go/internal/models"
	"task-scheduler-go/internal/scheduler"
	"task-scheduler-go/internal/storage"
)

func main() {
	// Initialize storage and scheduler
	store := storage.NewMemoryStorage()
	sched := scheduler.NewScheduler(store)

	// HTTP handlers.
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleCreateTask(sched, w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Start server
	port := ":8080"
	log.Printf("Scheduler server running on port %s", port)
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

	task, err := sched.SubmitTask(req.Type, req.Payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
